package main

import (
	"log"
	"os"

	"github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize RauthProvider
	config := &rauthprovider.Config{
		RauthAPIKey:   os.Getenv("RAUTH_API_KEY"),
		AppID:         os.Getenv("RAUTH_APP_ID"),
		WebhookSecret: os.Getenv("RAUTH_WEBHOOK_SECRET"),
	}

	if err := rauthprovider.Init(config); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	// Create Fiber app
	app := fiber.New()

	// Simple login endpoint
	app.Post("/login", func(c *fiber.Ctx) error {
		var req struct {
			SessionToken string `json:"session_token"`
			UserPhone    string `json:"user_phone"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// Verify session
		verified, err := rauthprovider.VerifySession(c.Context(), req.SessionToken, req.UserPhone)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Verification failed",
			})
		}

		if !verified {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid session",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"user":    req.UserPhone,
		})
	})

	// Protected endpoint with simple auth
	app.Get("/protected", func(c *fiber.Ctx) error {
		// Get token from header
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Missing token",
			})
		}

		// Remove "Bearer " prefix
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// Get phone from header
		phone := c.Get("X-User-Phone")
		if phone == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Missing phone",
			})
		}

		// Verify session
		verified, err := rauthprovider.VerifySession(c.Context(), token, phone)
		if err != nil || !verified {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid session",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Protected route accessed",
			"user":    phone,
		})
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		healthy, _ := rauthprovider.CheckAPIHealth(c.Context())
		return c.JSON(fiber.Map{
			"status": "healthy",
			"api":    healthy,
		})
	})

	log.Println("🚀 Simple Fiber server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
