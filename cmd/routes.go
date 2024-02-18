package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/leroysb/go_kubernetes/internal/api/handlers"
	"github.com/leroysb/go_kubernetes/internal/database"
)

func setupRoutes(app *fiber.App) {
	app.Use(cors.New())
	app.Use(logger.New())

	api := app.Group("/api/v1")
	api.Get("/status", StatusHandler)
	// api.Get("/stats", StatsHandler)

	// Product endpoints
	api.Get("/products", handlers.GetProducts)
	api.Post("/products", handlers.CreateProduct)
	api.Get("/products/:id", handlers.GetProduct)
	api.Put("/products/:id", handlers.UpdateProduct)
	api.Delete("/products/:id", handlers.DeleteProduct)

	// Customer endpoints
	// api.Post("/customers", CreateCustomer)
	// api.Get("/customers/me", GetCustomers)
	// api.Post("/customers/login", GetCustomer)
	// api.Post("/customers/logout", GetCustomer)

	// Order endpoints
	// api.Post("/orders", CreateOrder)
	// api.Get("/orders", GetOrders)
	// api.Put("/orders/:id", UpdateOrder)
	// api.Delete("/orders/:id", DeleteOrder)

	// 404 Handler
	app.Use(notFoundHandler)
}

func StatusHandler(c *fiber.Ctx) error {
	if database.CheckDBConnection() {
		return c.Status(200).JSON(fiber.Map{"Postgres": true})
	}
	return c.Status(500).JSON(fiber.Map{"Postgres": false})
}

func notFoundHandler(c *fiber.Ctx) error {
	return c.Status(404).JSON(fiber.Map{"error": "Not found"})
}
