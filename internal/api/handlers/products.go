package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leroysb/go_kubernetes/internal/database"
	"github.com/leroysb/go_kubernetes/internal/database/models"
)

type RequestBody struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Stock int    `json:"stock"`
}

// var db *gorm.DB

// CreateProduct creates a new product
func CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Check for required fields
	if product.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing product name"})
	}

	if product.Price == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing product price"})
	}

	if product.Stock == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing product stock"})
	}

	// Check if product already exists
	var existingProduct models.Product
	if err := database.DB.Db.Where("name = ?", product.Name).First(&existingProduct).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Product already exists"})
	}

	// Create the product in a goroutine
	go func() {
		if err := database.DB.Db.Create(&product).Error; err != nil {
			// Handle error in goroutine
			// fmt.Println("Error creating product:", err)
			c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
			return
		}
	}()

	return c.Status(200).JSON(product)
}

// GetProducts returns all products with pagination
func GetProducts(c *fiber.Ctx) error {
	page := c.Query("page")
	if page == "" {
		page = "1"
	}

	// Convert page to integer
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid page number"})
	}

	// Calculate offset and limit for pagination
	offset := (pageNum - 1) * 20
	limit := 20

	// Fetch products from database in a goroutine
	var products []models.Product
	done := make(chan bool)
	go func() {
		if err := database.DB.Db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
			done <- false
		} else {
			done <- true
		}
	}()

	select {
	case success := <-done:
		if success {
			return c.Status(200).JSON(products)
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
		}
	case <-time.After(5 * time.Second):
		return c.Status(500).JSON(fiber.Map{"error": "Timeout"})
	}
}

// GetProduct returns a single product
func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	product := new(models.Product)

	done := make(chan bool)
	go func() {
		if err := database.DB.Db.First(&product, id).Error; err != nil {
			done <- false
		} else {
			done <- true
		}
	}()

	select {
	case success := <-done:
		if success {
			return c.Status(200).JSON(product)
		} else {
			return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}
	case <-time.After(5 * time.Second):
		return c.Status(500).JSON(fiber.Map{"error": "Timeout"})
	}
}

// UpdateProduct updates a product
func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	product := new(models.Product)
	done := make(chan bool)

	go func() {
		database.DB.Db.First(&product, id)
		done <- true
	}()

	select {
	case <-done:
		if err := c.BodyParser(product); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		database.DB.Db.Save(&product)
		return c.JSON(product)
	case <-time.After(5 * time.Second):
		return c.Status(500).JSON(fiber.Map{"error": "Timeout"})
	}
}

// DeleteProduct deletes a product
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	product := new(models.Product)

	done := make(chan bool)
	go func() {
		if err := database.DB.Db.First(&product, id).Error; err != nil {
			done <- false
		} else {
			done <- true
		}
	}()

	select {
	case success := <-done:
		if success {
			database.DB.Db.Delete(&product)
			return c.Status(204).JSON(fiber.Map{"message": "Product deleted successfully"})
		} else {
			return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}
	case <-time.After(5 * time.Second):
		return c.Status(500).JSON(fiber.Map{"error": "Timeout"})
	}
}
