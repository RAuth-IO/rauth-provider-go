package main

import (
	"log"
	"os"

	"github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize the RauthProvider
	config := &rauthprovider.Config{
		RauthAPIKey:       os.Getenv("RAUTH_API_KEY"),
		AppID:             os.Getenv("RAUTH_APP_ID"),
		WebhookSecret:     os.Getenv("RAUTH_WEBHOOK_SECRET"),
		DefaultSessionTTL: 900,  // 15 minutes
		DefaultRevokedTTL: 3600, // 1 hour
	}

	if err := rauthprovider.Init(config); err != nil {
		log.Fatalf("Failed to initialize RauthProvider: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-User-Phone",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Health check endpoint
	app.Get("/health", healthHandler)

	// Webhook endpoint
	app.Post("/rauth/webhook", rauthWebhookHandler)

	// Session verification endpoint
	app.Post("/api/login", loginHandler)

	// Protected routes group
	protected := app.Group("/api", fiberAuthMiddleware())
	{
		protected.Get("/protected", protectedHandler)
		protected.Get("/user/profile", userProfileHandler)
		protected.Post("/logout", logoutHandler)
	}

	// Optional auth routes (works with or without authentication)
	optionalAuth := app.Group("/api/optional", fiberOptionalAuthMiddleware())
	{
		optionalAuth.Get("/public", publicHandler)
		optionalAuth.Get("/user-info", userInfoHandler)
	}

	log.Println("🚀 Fiber server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}

// Custom error handler for Fiber
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": message,
	})
}

// Health check handler
func healthHandler(c *fiber.Ctx) error {
	// Check API health
	healthy, err := rauthprovider.CheckAPIHealth(c.Context())
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unhealthy",
			"error":  err.Error(),
		})
	}

	if !healthy {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unhealthy",
			"error":  "Rauth API is not responding",
		})
	}

	return c.JSON(fiber.Map{
		"status": "healthy",
		"api":    "connected",
		"server": "fiber",
	})
}

// Webhook handler for Rauth events
func rauthWebhookHandler(c *fiber.Ctx) error {
	// For now, just return success
	// TODO: Implement proper webhook handling
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Webhook received",
	})
}

// Login handler for session verification
func loginHandler(c *fiber.Ctx) error {
	var req struct {
		SessionToken string `json:"session_token"`
		UserPhone    string `json:"user_phone"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Verify session
	verified, err := rauthprovider.VerifySession(c.Context(), req.SessionToken, req.UserPhone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Session verification failed",
			"details": err.Error(),
		})
	}

	if !verified {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Phone number not verified",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Session verified successfully",
		"user":    req.UserPhone,
	})
}

// Protected handler (requires authentication)
func protectedHandler(c *fiber.Ctx) error {
	// Get session info from context (set by middleware)
	sessionToken := c.Locals("session_token").(string)
	userPhone := c.Locals("user_phone").(string)

	return c.JSON(fiber.Map{
		"message": "Protected route accessed successfully",
		"user":    userPhone,
		"session": sessionToken,
		"server":  "fiber",
	})
}

// User profile handler
func userProfileHandler(c *fiber.Ctx) error {
	userPhone := c.Locals("user_phone").(string)

	return c.JSON(fiber.Map{
		"user": userPhone,
		"profile": fiber.Map{
			"phone":      userPhone,
			"verified":   true,
			"last_login": "2024-01-01T00:00:00Z",
		},
	})
}

// Logout handler
func logoutHandler(c *fiber.Ctx) error {
	// For now, just return success
	// TODO: Implement session revocation
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

// Public handler (optional auth)
func publicHandler(c *fiber.Ctx) error {
	// Check if user is authenticated
	_, hasSession := c.Locals("session_token").(string)
	userPhone, hasUser := c.Locals("user_phone").(string)

	response := fiber.Map{
		"message": "Public route accessed",
		"server":  "fiber",
	}

	if hasSession && hasUser {
		response["authenticated"] = true
		response["user"] = userPhone
	} else {
		response["authenticated"] = false
	}

	return c.JSON(response)
}

// User info handler (optional auth)
func userInfoHandler(c *fiber.Ctx) error {
	userPhone, hasUser := c.Locals("user_phone").(string)

	if !hasUser {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Authentication required for user info",
		})
	}

	return c.JSON(fiber.Map{
		"user": userPhone,
		"info": fiber.Map{
			"phone":         userPhone,
			"authenticated": true,
		},
	})
}

// Fiber authentication middleware
func fiberAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract session token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Missing authorization header",
			})
		}

		// Check if it's a Bearer token
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid authorization header format",
			})
		}

		sessionToken := authHeader[7:] // Remove "Bearer " prefix

		// Extract user phone from header or query parameter
		userPhone := c.Get("X-User-Phone")
		if userPhone == "" {
			userPhone = c.Query("user_phone")
		}

		if userPhone == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Missing user phone",
			})
		}

		// Verify session
		verified, err := rauthprovider.VerifySession(c.Context(), sessionToken, userPhone)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Session verification failed",
			})
		}

		if !verified {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid session",
			})
		}

		// Check if session is revoked
		isRevoked, err := rauthprovider.IsSessionRevoked(c.Context(), sessionToken)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Session status check failed",
			})
		}

		if isRevoked {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Session revoked",
			})
		}

		// Add session info to context
		c.Locals("session_token", sessionToken)
		c.Locals("user_phone", userPhone)

		return c.Next()
	}
}

// Fiber optional authentication middleware
func fiberOptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract session token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// No auth header, continue without session
			return c.Next()
		}

		// Check if it's a Bearer token
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			// Invalid format, continue without session
			return c.Next()
		}

		sessionToken := authHeader[7:] // Remove "Bearer " prefix

		// Extract user phone from header or query parameter
		userPhone := c.Get("X-User-Phone")
		if userPhone == "" {
			userPhone = c.Query("user_phone")
		}

		if userPhone == "" {
			// No user phone, continue without session
			return c.Next()
		}

		// Try to verify session
		verified, err := rauthprovider.VerifySession(c.Context(), sessionToken, userPhone)
		if err != nil || !verified {
			// Verification failed, continue without session
			return c.Next()
		}

		// Check if session is revoked
		isRevoked, err := rauthprovider.IsSessionRevoked(c.Context(), sessionToken)
		if err != nil || isRevoked {
			// Session revoked or error, continue without session
			return c.Next()
		}

		// Add session info to context
		c.Locals("session_token", sessionToken)
		c.Locals("user_phone", userPhone)

		return c.Next()
	}
}
