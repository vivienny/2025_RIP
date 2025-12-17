package main

import (
	"decode/internal/app/ds"
	"decode/internal/app/dsn"
	"log"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt" // ДОБАВИТЬ этот импорт
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Add role column if not exists
	if !db.Migrator().HasColumn(&ds.Avius{}, "role") {
		db.Exec("ALTER TABLE avius ADD COLUMN role VARCHAR(20) DEFAULT 'user'")
		log.Println("✅ Добавлено поле role в таблицу avius")
	}

	// Update existing users
	db.Exec("UPDATE avius SET role = 'moderator' WHERE is_moderator = true")
	db.Exec("UPDATE avius SET role = 'user' WHERE role IS NULL OR role = ''")

	log.Println("✅ Роли пользователей обновлены")

	// Run full migration
	err = db.AutoMigrate(
		&ds.ASGARService{},
		&ds.Avius{},
		&ds.FlightServ{},
		&ds.Subjserv{},
	)
	if err != nil {
		panic("cant migrate db")
	}

	// ========== СОЗДАЕМ ТЕСТОВЫХ ПОЛЬЗОВАТЕЛЕЙ ==========
	testUsers := []ds.Avius{
		{
			Login:    "Meow",
			Password: "meow123",
			Role:     "user",
		},
		{
			Login:       "moderator1",
			Password:    "moderator123",
			Role:        "moderator",
			IsModerator: true,
		},
	}

	for _, user := range testUsers {
		// Проверяем, существует ли уже пользователь
		var existingUser ds.Avius
		result := db.Where("login = ?", user.Login).First(&existingUser)

		if result.Error == gorm.ErrRecordNotFound {
			// Пользователя нет - создаем нового с хешированным паролем
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("⚠️ Ошибка хеширования пароля для %s: %v", user.Login, err)
				continue
			}
			user.Password = string(hashedPassword)

			if err := db.Create(&user).Error; err != nil {
				log.Printf("⚠️ Ошибка создания пользователя %s: %v", user.Login, err)
			} else {
				log.Printf("✅ Создан пользователь: %s (роль: %s)", user.Login, user.Role)
			}
		} else {
			// Пользователь уже существует - обновляем роль если нужно
			if existingUser.Role != user.Role {
				db.Model(&existingUser).Update("role", user.Role)
				log.Printf("🔄 Обновлена роль пользователя %s на %s", user.Login, user.Role)
			} else {
				log.Printf("ℹ️ Пользователь %s уже существует (роль: %s)", user.Login, existingUser.Role)
			}
		}
	}

	log.Println("🎉 Миграция ролей и тестовых пользователей завершена успешно!")
}
