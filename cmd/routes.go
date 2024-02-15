package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	app.Use(cors.New())
	app.Use(logger.New())

	api := app.Group("/api")
	api.Get("/status", StatusHandler)
	// api.Post("/customers", CreateCustomer)
	// api.Get("/customers", GetCustomers)
	// api.Get("/customers/:id", GetCustomer)
	// api.Post("/products", CreateProducts)
	// api.Get("/products", GetProducts)
	// api.Get("/products/:id", GetProduct)
	// api.Put("/products/:id", UpdateProduct)
	// api.Delete("/products/:id", DeleteProduct)
	// api.Post("/orders", CreateOrder)
	// api.Get("/orders", GetOrders)
	// api.Get("/orders/:id", GetOrder)
	// api.Put("/orders/:id", UpdateOrder)
	// api.Delete("/orders/:id", DeleteOrder)
}

func StatusHandler(c *fiber.Ctx) error {
	if db.DB != nil {
		return c.Status(200).JSON(fiber.Map{"Postgres": true})
	}
	return c.Status(500).JSON(fiber.Map{"Postgres": false}) 
}
