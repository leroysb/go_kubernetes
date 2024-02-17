package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/leroysb/go_kubernetes/internal/database"
)

func main() {
	// Connect to database
	database.ConnectDB()

	// Initialize Fiber app
	app := fiber.New()

	// Define routes
	setupRoutes(app)

	// Start server
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("API_PORT")

	if port == "" {
		port = "8080"
	}

	err = app.Listen(":" + port)

	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
