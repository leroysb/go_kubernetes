package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/leroysb/go_kubernetes/internal/api/auth"
	"github.com/leroysb/go_kubernetes/internal/api/handlers"
	"github.com/leroysb/go_kubernetes/internal/database"
)

func setupRoutes(app *fiber.App) {
	// Middleware
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded")
		},
	}))
	app.Use(cors.New())
	app.Use(logger.New())
	// app.Use(auth.AuthMiddleware())

	api := app.Group("/api/v1")
	api.Get("/status", StatusHandler)
	// api.Get("/stats", StatsHandler)

	// Public Product endpoints
	api.Get("/products", handlers.GetProducts)
	api.Post("/products", handlers.CreateProduct)
	api.Get("/products/:id", handlers.GetProduct)
	api.Put("/products/:id", handlers.UpdateProduct)
	api.Delete("/products/:id", handlers.DeleteProduct)

	// Customer endpoints
	api.Post("/customers", handlers.CreateCustomer) // Public Endpoint for user registration

	api.Get("/customers/login", handlers.Login)                           // Public Endpoint for user authentication
	api.Get("/customers/me", auth.AuthMiddleware(), handlers.GetCustomer) // Endpoint to retrieve authorized user information
	api.Post("/customers/logout", auth.AuthMiddleware(), handlers.Logout)
	api.Post("/customers/cart", handlers.CreateCart)
	api.Get("/customers/cart", handlers.GetCart)
	api.Put("/customers/cart/:id", handlers.UpdateCart)
	api.Delete("/customers/cart/:id", handlers.DeleteCart)
	api.Post("/customers/orders/:id", handlers.CreateOrder)

	// Order endpoints
	api.Get("/orders", handlers.GetOrders)

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
