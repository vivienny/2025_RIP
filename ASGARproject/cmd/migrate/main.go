package main

import (
	"decode/internal/app/ds"
	"decode/internal/app/dsn"

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

	// Migrate the schema - добавляем новые модели
	err = db.AutoMigrate(
		&ds.ASGARService{}, // уже была
		&ds.Avius{},        // новая
		&ds.Miniplane{},    // новая
		&ds.FlightServ{},   // новая
		&ds.Subjserv{},     // новая
	)
	if err != nil {
		panic("cant migrate db")
	}
}
