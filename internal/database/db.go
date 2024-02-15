package database

import (
	"fmt"
	"log"
	"os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"gorm.io/gorm/logger"
	"/internal/models"
)

type DBinstance struct {
	DB *gorm.DB
}

var DB *DBinstance

// InitializeDB initializes the database connection and performs migrations
func ConnectDB() {
	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
    if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		// os.Exit(11)
        return nil, err
    }

	log.Println("Connected to database")
	db.Logger = logger.Default.LogMode(logger.Info)

    // Perform auto-migration
	log.Println("Performing auto-migration")
	
    err = MigrateDB()
    if err != nil {
		log.Fatalf("Error performing auto-migration: %v", err)
        return nil, err
    }

	return db, nil
}

// MigrateDB performs auto-migration for all models
func MigrateDB() error {
    return db.AutoMigrate(&models.Order{}, &models.Product{}, &models.Customer{})
}

// GetDB returns the initialized database instance
func GetDB() *gorm.DB {
	return db
}
