# 🚀 RauthProvider Go Library v1.0.0

A lightweight, plug-and-play Go library for phone number authentication using the Rauth.io reverse verification flow via WhatsApp or SMS.

## ✨ Features

### 🔐 Core Authentication
- **Reverse Authentication** - Authenticate users via WhatsApp or SMS without sending OTPs
- **Session Management** - Track sessions, verify tokens, and revoke access automatically
- **Webhook Support** - Listen for number verification and session revocation in real-time
- **Signature Verification** - HMAC-SHA256 signature verification for webhooks

### 🏗️ Architecture & Design
- **Clean Architecture** - Built with clean architecture principles and Go best practices
- **Singleton Pattern** - Thread-safe singleton provider using `sync.Once`
- **Dependency Injection** - Interface-based dependencies for testability
- **Repository Pattern** - Abstract data access layer
- **Middleware Pattern** - Chainable HTTP middleware

### ⚡ Performance & Security
- **In-Memory Storage** - Sessions stored in memory for fast access
- **Automatic Cleanup** - Expired sessions cleaned up every 5 minutes
- **Thread-Safe** - All operations thread-safe using read-write mutexes
- **Connection Pooling** - HTTP client uses connection pooling for API calls
- **Context Support** - All operations support context for cancellation and timeouts

### 🧩 Developer Experience
- **Plug-and-Play API** - Simple, developer-friendly API surface
- **HTTP Middleware** - Drop-in HTTP middleware integration
- **Comprehensive Examples** - Basic HTTP, Fiber, and Gin framework examples
- **Error Handling** - Detailed error types and custom error handling
- **Documentation** - Complete API reference and usage examples

## 📦 Installation

```bash
go get github.com/RAuth-IO/rauth-provider-go@v1.0.0
```

## 🚀 Quick Start

```go
package main

import (
    "log"
    "os"
    "github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
)

func main() {
    config := &rauthprovider.Config{
        RauthAPIKey:   os.Getenv("RAUTH_API_KEY"),
        AppID:         os.Getenv("RAUTH_APP_ID"),
        WebhookSecret: os.Getenv("RAUTH_WEBHOOK_SECRET"),
    }

    if err := rauthprovider.Init(config); err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }

    // Verify session
    verified, err := rauthprovider.VerifySession(ctx, "session-token", "+1234567890")
    if err != nil {
        log.Printf("Verification failed: %v", err)
        return
    }

    if verified {
        log.Println("Session is valid")
    }
}
```

## 🔧 Framework Integration

### Fiber Framework
```go
// Simple integration
app.Post("/login", func(c *fiber.Ctx) error {
    var req struct {
        SessionToken string `json:"session_token"`
        UserPhone    string `json:"user_phone"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    verified, err := rauthprovider.VerifySession(c.Context(), req.SessionToken, req.UserPhone)
    if err != nil || !verified {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
    }

    return c.JSON(fiber.Map{"success": true, "user": req.UserPhone})
})
```

### Gin Framework
```go
// Protected routes with middleware
protected := r.Group("/api")
protected.Use(ginAuthMiddleware())
{
    protected.GET("/protected", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Protected route"})
    })
}
```

## 📊 API Reference

### Core Functions
- `rauthprovider.Init(config *rauthprovider.Config) error` - Initialize provider
- `rauthprovider.VerifySession(ctx, sessionToken, userPhone) (bool, error)` - Verify session
- `rauthprovider.IsSessionRevoked(ctx, sessionToken) (bool, error)` - Check revocation
- `rauthprovider.CheckAPIHealth(ctx) (bool, error)` - Check API health
- `rauthprovider.WebhookHandler() http.HandlerFunc` - Webhook handler

### Middleware Functions
- `middleware.AuthMiddleware() func(http.Handler) http.Handler` - Required auth
- `middleware.OptionalAuthMiddleware() func(http.Handler) http.Handler` - Optional auth
- `middleware.GetSessionToken(ctx) (string, bool)` - Extract session token
- `middleware.GetUserPhone(ctx) (string, bool)` - Extract user phone

## 📁 Project Structure

```
├── pkg/                    # Public API
│   ├── rauthprovider/     # Main provider package
│   └── middleware/        # HTTP middleware
├── internal/              # Private application code
│   ├── domain/           # Business entities and interfaces
│   ├── usecase/          # Application business rules
│   ├── infrastructure/   # External interfaces (API, storage)
│   └── delivery/         # HTTP handlers, webhooks
└── examples/             # Usage examples
    ├── basic/            # Basic HTTP server example
    └── fiber/            # Fiber framework examples
        ├── simple/       # Simple integration
        └── advanced/     # Advanced with middleware
```

## 🔒 Security Features

- **HMAC-SHA256 Signature Verification** - Webhook signatures verified
- **Session Validation** - Sessions validated against phone numbers
- **Revocation Tracking** - Revoked sessions tracked and checked
- **TTL Management** - Automatic expiration of sessions
- **Input Validation** - All inputs validated before processing

## 🧪 Testing

```go
func TestSessionVerification(t *testing.T) {
    config := &rauthprovider.Config{
        RauthAPIKey:   "test-key",
        AppID:         "test-app",
        WebhookSecret: "test-secret",
    }
    
    if err := rauthprovider.Init(config); err != nil {
        t.Fatalf("Failed to initialize: %v", err)
    }

    verified, err := rauthprovider.VerifySession(context.Background(), "test-token", "+1234567890")
    if err != nil {
        t.Errorf("Verification failed: %v", err)
    }

    t.Logf("Verification result: %v", verified)
}
```

## 📈 Performance Considerations

- **In-Memory Storage** - Sessions stored in memory for fast access
- **Automatic Cleanup** - Expired sessions cleaned up every 5 minutes
- **Thread-Safe** - All operations thread-safe using read-write mutexes
- **Connection Pooling** - HTTP client uses connection pooling for API calls
- **Context Support** - All operations support context for cancellation and timeouts

## 🔄 Breaking Changes

- **Updated Config Type** - Now uses `rauthprovider.Config` instead of internal `domain.Config`
- **Import Paths** - All examples updated to use correct import paths
- **Public API** - Clean separation between public and internal APIs

## 📚 Documentation

- **Complete API Reference** - Available in README.md
- **Framework Examples** - Fiber, Gin, and basic HTTP examples
- **Test Scripts** - Automated testing scripts for examples
- **Error Handling** - Comprehensive error types and handling

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details.

## 🆘 Support

- **GitHub Issues**: [Create an issue](https://github.com/RAuth-IO/rauth-provider-go/issues)
- **Documentation**: [Rauth.io Documentation](https://docs.rauth.io)
- **Email**: support@rauth.io

---

**Ready for production use!** 🚀

This release provides a complete, production-ready Go library for Rauth authentication with comprehensive examples, documentation, and framework integrations. 