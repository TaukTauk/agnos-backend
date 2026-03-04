package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"agnos-backend/internal/model"
)

// Load reads the .env file into environment variables
func Load() {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, reading from environment")
	}
}

// ConnectDB creates a connection to PostgreSQL using GORM
func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

// Migrate runs GORM AutoMigrate for all models
func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Hospital{},
		&model.Staff{},
		&model.Patient{},
	)
	log.Println("Database migrated successfully")
}