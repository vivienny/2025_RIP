package main

import (
	"decode/internal/app/ds"
	"decode/internal/app/dsn"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 1. Сначала добавляем flight_serv_id как NULLABLE
	if !db.Migrator().HasColumn(&ds.Subjserv{}, "flight_serv_id") {
		db.Exec("ALTER TABLE subjservs ADD COLUMN flight_serv_id BIGINT")
		log.Println("✅ Добавлен flight_serv_id")
	}

	// 2. Переносим данные из miniplane_id в flight_serv_id
	db.Exec("UPDATE subjservs SET flight_serv_id = miniplane_id WHERE flight_serv_id IS NULL")
	log.Println("✅ Данные перенесены из miniplane_id в flight_serv_id")

	// 4. Теперь делаем flight_serv_id NOT NULL
	db.Exec("ALTER TABLE subjservs ALTER COLUMN flight_serv_id SET NOT NULL")
	log.Println("✅ flight_serv_id установлен как NOT NULL")

	// 5. Удаляем старый столбец miniplane_id
	if db.Migrator().HasColumn(&ds.Subjserv{}, "miniplane_id") {
		db.Exec("ALTER TABLE subjservs DROP COLUMN miniplane_id")
		log.Println("✅ Столбец miniplane_id удален")
	}

	// 6. Добавляем статус в flight_servs если нет
	if !db.Migrator().HasColumn(&ds.FlightServ{}, "status") {
		db.Exec("ALTER TABLE flight_servs ADD COLUMN status VARCHAR(15) NOT NULL DEFAULT 'cart'")
		log.Println("✅ Добавлен статус в flight_servs")
	}

	// 7. Создаем/обновляем остальные таблицы
	err = db.AutoMigrate(
		&ds.ASGARService{}, // услуги
		&ds.Avius{},        // пользователи
		&ds.FlightServ{},   // заявки
		&ds.Subjserv{},     // позиции заказа
	)
	if err != nil {
		panic("cant migrate db")
	}

	log.Println("🎉 Миграция завершена успешно!")
}
