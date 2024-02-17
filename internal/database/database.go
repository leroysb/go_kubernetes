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

	// Load environment variables
	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
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
		// log.Fatalf("Failed to connect to database: %v", err)
		// err = fmt.Errorf("failed to connect to database: %v", err)
		// log.Println(err)
		log.Printf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)

	// Perform auto-migration
	log.Println("Performing auto-migration")
	if err := MigrateDB(); err != nil {
		// log.Fatalf("Error performing auto-migration: %v", err)
		// err = fmt.Errorf("error performing auto-migration: %v", err)
		// log.Println(err)
		log.Printf("Error performing auto-migration: %v", err)
	}

	log.Println("Database migration successful")
}

// MigrateDB performs auto-migration for all models
func MigrateDB() error {
	return db.AutoMigrate(&models.Product{})
}

// GetDB returns the initialized database instance
func GetDB() *gorm.DB {
	return db
}

func CheckDBConnection() bool {
	if db == nil {
		log.Println("Database connection is nil")
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting underlying database connection: %v", err)
		return false
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return false
	}

	log.Println("Database connection is alive")
	return true
}
