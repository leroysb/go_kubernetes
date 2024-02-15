package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/leroysb/go_kubernetes/internal/database"
)

// var db *sql.DB
var db *gorm.DB

func main() {
	// Connect to database
	db = database.GetDB()
	
	// Initialize Fiber app
	app := fiber.New()

	// Define routes
	setupRoutes(app)

	// Start server
	port := os.Getenv("API_PORT")
	
	if port == "" {
		port = "8080"
	}

	err = app.Listen(":" + port)

	if err != nil {
		log.Fatal("Error starting server:", err)
	} else {
		log.Println("Server is running on port", port)
		fmt.Println("Server is running on port", port)
	}
}
