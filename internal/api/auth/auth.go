package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TokenInfo struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope"`
	ClientID  string `json:"client_id"`
	Sub       string `json:"sub"`
	Exp       int    `json:"exp"`
	Iat       int    `json:"iat"`
	Nbf       int    `json:"nbf"`
	Aud       []any  `json:"aud"`
	Iss       string `json:"iss"`
	TokenType string `json:"token_type"`
	TokenUse  string `json:"token_use"`
}

var requiredScope = os.Getenv("requiredScope")
var hydraAdminUrl = os.Getenv("hydraAdminUrl")

// introspectToken sends a request to the Hydra introspection endpoint to validate the access token
func introspectToken(accessToken string) (*TokenInfo, error) {
	// Prepare the form data
	formData := url.Values{}
	formData.Set("token", accessToken)
	formData.Set("scope", requiredScope)

	// Create the HTTP request
	req, err := http.NewRequest("POST", hydraAdminUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("Error creating introspection request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Send the request
	introspectResp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending introspection request:", err)
		return nil, err
	}
	defer introspectResp.Body.Close()

	// Check the response status code
	if introspectResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("introspection failed with status code %d", introspectResp.StatusCode)
	}

	// Parse introspection response
	var tokenInfo TokenInfo
	if err := json.NewDecoder(introspectResp.Body).Decode(&tokenInfo); err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

// checkTokenActive checks if the introspected token is active
func hasScope(tokenScope string, requiredScope string) bool {
	scopes := strings.Split(tokenScope, " ")
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}
	return false
}

// AuthMiddleware is a middleware function to validate access token using Hydra introspection endpoint
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the access token from the request headers
		authHeader := c.Get("Authorization")

		// Check if Authorization header is missing or does not start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Extract the token by stripping the "Bearer " prefix
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Introspect the token
		tokenInfo, err := introspectToken(accessToken)
		if err != nil || !tokenInfo.Active {
			fmt.Println("Error during token introspection:", err)
			c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
			return err
		}

		// Check if the token is active
		if !hasScope(tokenInfo.Scope, requiredScope) {
			// If token is not active or not present, return unauthorized error
			fmt.Printf("Insufficient scope: %v\n", err)
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Insufficient scope"})
			// return err
		}

		// Proceed to the next middleware or handler if token is valid
		return c.Next()
	}
}
