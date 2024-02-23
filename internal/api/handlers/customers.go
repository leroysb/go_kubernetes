package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leroysb/go_kubernetes/internal/database"
	"github.com/leroysb/go_kubernetes/internal/database/models"
	"github.com/leroysb/go_kubernetes/internal/sms"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type HydraClientRequest struct {
	ClientName        string   `json:"client_name"`
	ClientSecret      string   `json:"client_secret"`
	GrantTypes        []string `json:"grant_types"`
	Scope             string   `json:"scope"`
	TokenEndpointAuth string   `json:"token_endpoint_auth_method"`
}

type HydraClientResponse struct {
	ClientID string `json:"client_id"`
}

type Cart struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

// checkPasswordHash compares a hashed password with its possible plaintext equivalent
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// hashPassword generates a hash of the password using bcrypt
func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

// getUserByEmail retrieves a user from the database by email
func getUserByPhone(phone string) (*models.Customer, error) {
	var user models.Customer
	if err := database.DB.Db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return nil if user not found
			return nil, nil
		}
		// Return error for other database errors
		return nil, err
	}
	return &user, nil
}

// CreateCustomer creates a new customer
func CreateCustomer(c *fiber.Ctx) error {
	customer := new(models.Customer)

	// Error check fields
	if err := c.BodyParser(customer); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			if strings.Contains(err.Error(), "name") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing name of type string"})
			}
			if strings.Contains(err.Error(), "phone") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing phone number of type string"})
			}
			if strings.Contains(err.Error(), "password") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing password of type string"})
			}
		}
		return c.Status(400).SendString(err.Error())
	}

	if customer.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing name"})
	}

	if customer.Phone == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing phone number"})
	}

	if customer.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing password"})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	customer.Password = string(hashedPassword)

	// Check if customer already exists
	var existingCustomer models.Customer
	if err := database.DB.Db.Where("phone = ?", customer.Phone).First(&existingCustomer).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Customer already exists"})
	}

	// Create the customer in a goroutine
	go func() {
		if err := database.DB.Db.Create(&customer).Error; err != nil {
			// Handle error in goroutine
			// fmt.Println("Error creating customer:", err)
			c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
			return
		}
	}()

	go func() {
		sms.SendSMS(customer.Phone, "Welcome to our platform")
	}()

	return c.Status(200).JSON(customer)
}

func Login(c *fiber.Ctx) error {
	var loginData struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&loginData); err != nil {
		return err
	}

	// Retrieve user from the database
	user, err := getUserByPhone(loginData.Phone)
	if err != nil {
		return err
	}

	// Check if the user exists and the password matches
	if user == nil || !checkPasswordHash(loginData.Password, user.Password) {
		// Return unauthorized error if user does not exist or password does not match
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Create request body for Hydra client creation
	clientRequest := HydraClientRequest{
		ClientName:        user.Name,
		ClientSecret:      hashPassword(loginData.Password),
		GrantTypes:        []string{"client_credentials"},
		Scope:             os.Getenv("hydraScope"),
		TokenEndpointAuth: "client_secret_post",
	}

	// Convert struct to JSON
	clientRequestBody, err := json.Marshal(clientRequest)
	if err != nil {
		return err
	}

	// Send POST request to Hydra client creation endpoint
	resp, err := http.Post(os.Getenv("hydraClientUrl"), "application/json", bytes.NewBuffer(clientRequestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response
	var clientResponse HydraClientResponse
	if err := json.NewDecoder(resp.Body).Decode(&clientResponse); err != nil {
		return err
	}

	// Return the client_id as the response
	// return c.JSON(clientResponse)

	// Create request body for Hydra token creation
	tokenRequest := struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Scope        string `json:"scope"`
	}{
		GrantType:    "client_credentials",
		ClientID:     clientResponse.ClientID,
		ClientSecret: clientRequest.ClientSecret, // Use the hashed password as client secret
		Scope:        os.Getenv("hydraScope"),
	}

	// Convert struct to JSON
	tokenRequestBody, err := json.Marshal(tokenRequest)
	if err != nil {
		return err
	}

	// Send POST request to Hydra token creation endpoint
	resp, err = http.Post(os.Getenv("hydraTokenUrl"), "application/json", bytes.NewBuffer(tokenRequestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response
	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}

	// Extract access token from response
	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return errors.New("access token not found in response")
	}
	// Set the access token in the response headers
	c.Set("Authorization", "Bearer "+accessToken)

	// Return a success message
	return c.JSON(fiber.Map{"message": "Login successful"})
}

// GetCustomer retrieves information of currently authorized user
func GetCustomer(c *fiber.Ctx) error {
	// Retrieve user information from the context
	user := c.Locals("user").(*models.Customer)

	// Return user information
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	// customer := new(models.Customer)
	return nil
}

func CreateCart(c *fiber.Ctx) error {
	order := new(models.Order)

	// Set the customer_id from the authorized user
	user := c.Locals("user").(*models.Customer)
	order.CustomerID = user.ID

	// Error check fields
	if err := c.BodyParser(order); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			if strings.Contains(err.Error(), "product_id") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing product_id of type int"})
			}
			if strings.Contains(err.Error(), "quantity") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing quantity of type int"})
			}
		}
		// return c.Status(400).SendString(err.Error())
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if order.ProductID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing product_id"})
	}

	if order.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing quantity"})
	}

	// Retrieve product from the database
	var product models.Product
	if err := database.DB.Db.Where("id = ?", order.ProductID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(400).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// check if product is available
	if product.Stock < order.Quantity {
		return c.Status(400).JSON(fiber.Map{"error": "Product not available"})
	}

	// Set the product_id and amount
	order.ProductID = product.ID
	order.Amount = order.Quantity * product.Price

	// Set the time and status
	order.Status = "cart"

	// Create the order in a goroutine
	go func() {
		if err := database.DB.Db.Create(&order).Error; err != nil {
			// Handle error in goroutine
			// fmt.Println("Error creating order:", err)
			c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
			return
		}
	}()

	return c.Status(200).JSON(order)
}

func GetCart(c *fiber.Ctx) error {
	// Retrieve user information from the context
	user := c.Locals("user").(*models.Customer)

	// Retrieve cart items from the database
	var cartItems []Cart
	if err := database.DB.Db.Table("orders").Select("product_id, quantity").Where("customer_id = ? AND status = ?", user.ID, "cart").Scan(&cartItems).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Return cart items
	return c.JSON(cartItems)
}

func UpdateCart(c *fiber.Ctx) error {
	// Retrieve user information from the context
	user := c.Locals("user").(*models.Customer)

	// Retrieve product_id and quantity from the request
	var cartItem Cart
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Retrieve the order from the database
	var order models.Order
	if err := database.DB.Db.Where("id = ? AND customer_id = ? AND status = ?", c.Params("id"), user.ID, "cart").First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(400).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Update the order in the database
	if err := database.DB.Db.Model(&order).Updates(models.Order{Quantity: cartItem.Quantity}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(order)
}

func DeleteCart(c *fiber.Ctx) error {
	// Retrieve user information from the context
	user := c.Locals("user").(*models.Customer)

	// Retrieve the order from the database
	var order models.Order
	if err := database.DB.Db.Where("id = ? AND customer_id = ? AND status = ?", c.Params("id"), user.ID, "cart").First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(400).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Delete the order from the database
	if err := database.DB.Db.Delete(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Order deleted"})
}

func CreateOrder(c *fiber.Ctx) error {
	// Retrieve user information from the context
	user := c.Locals("user").(*models.Customer)

	order := new(models.Order)

	// Error check fields
	if err := c.BodyParser(order); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			if strings.Contains(err.Error(), "product_id") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing product_id of type int"})
			}
		}
		// return c.Status(400).SendString(err.Error())
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if order.ProductID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing product_id"})
	}

	// Retrieve product from the database
	var product models.Product
	if err := database.DB.Db.Where("id = ?", order.ProductID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(400).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// check if product is available
	if product.Stock < order.Quantity {
		return c.Status(400).JSON(fiber.Map{"error": "Product not available"})
	}

	// Set the product_id and amount
	order.ProductID = product.ID
	order.Amount = order.Quantity * product.Price

	// Set the time and status
	order.Time = time.Now().Format("2006-01-02 15:04:05")
	order.Status = "ordered"

	// Create the order in a goroutine
	go func() {
		if err := database.DB.Db.Create(&order).Error; err != nil {
			// Handle error in goroutine
			fmt.Println("Error creating order:", err)
			c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
			return
		}
	}()

	go func() {
		sms.SendSMS(user.Phone, "Order successful")
	}()

	// reduce the stock of the product
	product.Stock -= order.Quantity
	if err := database.DB.Db.Model(&product).Updates(models.Product{Stock: product.Stock}).Error; err != nil {
		fmt.Println("Error updating product stock:", err)
	}

	return c.Status(200).JSON(order)
}
