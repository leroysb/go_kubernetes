package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/leroysb/go_kubernetes/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// InitializeDB initializes the database connection and performs migrations
func ConnectDB() {
	var err error

	err = godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
		return
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}

	log.Println("Connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)

	// Perform auto-migration
	log.Println("Performing auto-migration")
	db.AutoMigrate(&models.Product{}, &models.Customer{}, &models.Order{})

	log.Println("Database migration successful")

}

func CheckDBConnection() bool {
	sqlDB, _ := db.DB()

	err := sqlDB.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return false
	}

	log.Println("Database connection is alive")
	return true
}
