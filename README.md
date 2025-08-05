# RauthProvider (Go)

A lightweight, plug-and-play Go library for phone number authentication using the Rauth.io reverse verification flow via WhatsApp or SMS. It handles everything from session creation to revocation — with real-time webhook updates and HTTP middleware support, all with minimal setup.

---

## ✅ Features

📲 **Reverse Authentication** – Authenticate users via WhatsApp or SMS without sending OTPs

🔐 **Session Management** – Track sessions, verify tokens, and revoke access automatically

📡 **Webhook Support** – Listen for number verification and session revocation in real-time

🧩 **Plug-and-Play API** – Simple, developer-friendly API surface

⚡ **HTTP Middleware** – Drop-in HTTP middleware integration

🛡️ **Secure by Design** – Webhook secret verification and session validation

🧠 **Smart Caching** – In-memory session tracking with fallback to API

🔗 **Rauth API Ready** – Built to connect seamlessly with the Rauth.io platform

🟪 **Clean Architecture** – Built with clean architecture principles and Go best practices

🔄 **Node.js Compatible** – Webhook authentication matches Node.js implementation

---

## Recent Updates

### ✅ **Webhook Authentication Simplified**
- **Changed from**: HMAC-SHA256 signature verification
- **Changed to**: Simple webhook secret comparison (Node.js style)
- **Header**: `x-webhook-secret` (matches Node.js implementation)
- **Benefit**: Easier testing and Node.js compatibility

### ✅ **API Endpoints Fixed**
- **Base URL**: Updated to `https://api.rauth.io/session`
- **Session Endpoint**: Changed from `/verify-session` to `/status`
- **Headers**: Added `X-App-ID` and browser-like headers
- **Cloudflare Protection**: Added retry logic and proper headers

### ✅ **Local Store Only**
- **`IsSessionRevoked()`**: Only checks local memory store
- **No API calls**: Never makes network requests for revocation checks
- **Fast Performance**: Instant response from local cache

---

## Installation

```bash
go get github.com/RAuth-IO/rauth-provider-go
```

---

## Quick Start

### Basic Usage
```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
)

func main() {
    // Initialize the provider
    config := &rauthprovider.Config{
        RauthAPIKey:      os.Getenv("RAUTH_API_KEY"),
        AppID:            os.Getenv("RAUTH_APP_ID"),
        WebhookSecret:    os.Getenv("RAUTH_WEBHOOK_SECRET"),
        DefaultSessionTTL: 900,  // 15 minutes
        DefaultRevokedTTL: 3600, // 1 hour
    }

    if err := rauthprovider.Init(config); err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }

    // Set up webhook handler
    http.HandleFunc("/rauth/webhook", rauthprovider.WebhookHandler())

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Session Verification
```go
// Verify a session
verified, err := rauthprovider.VerifySession(ctx, "session-token", "+1234567890")
if err != nil {
    log.Printf("Verification failed: %v", err)
    return
}

if verified {
    log.Println("Session is valid")
} else {
    log.Println("Session is invalid")
}
```

### Check Session Revocation
```go
// Check if session is revoked
isRevoked, err := rauthprovider.IsSessionRevoked(ctx, "session-token")
if err != nil {
    log.Printf("Check failed: %v", err)
    return
}

if isRevoked {
    log.Println("Session has been revoked")
} else {
    log.Println("Session is active")
}
```

---

## Usage Examples

### HTTP Server with Middleware
```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/RAuth-IO/rauth-provider-go/pkg/middleware"
    "github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
)

func main() {
    // Initialize provider
    config := &rauthprovider.Config{
        RauthAPIKey:      os.Getenv("RAUTH_API_KEY"),
        AppID:            os.Getenv("RAUTH_APP_ID"),
        WebhookSecret:    os.Getenv("RAUTH_WEBHOOK_SECRET"),
    }

    if err := rauthprovider.Init(config); err != nil {
        log.Fatal(err)
    }

    // Create server
    mux := http.NewServeMux()

    // Webhook endpoint
    mux.HandleFunc("/rauth/webhook", rauthprovider.WebhookHandler())

    // Session verification endpoint
    mux.HandleFunc("/api/login", loginHandler)

    // Protected routes with authentication middleware
    protectedMux := http.NewServeMux()
    protectedMux.HandleFunc("/api/protected", protectedHandler)
    
    authMiddleware := middleware.AuthMiddleware()
    protectedHandler := authMiddleware(protectedMux)
    mux.Handle("/api/", http.StripPrefix("/api", protectedHandler))

    log.Fatal(http.ListenAndServe(":8080", mux))
}

// Test webhook with proper authentication
func testWebhook() {
    payload := `{"event":"session_revoked","session_token":"test_token"}`
    req, _ := http.NewRequest("POST", "http://localhost:8080/rauth/webhook", strings.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-webhook-secret", "your-webhook-secret") // Node.js style authentication
    
    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        SessionToken string `json:"session_token"`
        UserPhone    string `json:"user_phone"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    verified, err := rauthprovider.VerifySession(r.Context(), req.SessionToken, req.UserPhone)
    if err != nil || !verified {
        middleware.WriteError(w, http.StatusUnauthorized, "invalid_session", "Phone number not verified")
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Session verified",
    })
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    sessionToken, _ := middleware.GetSessionToken(r.Context())
    userPhone, _ := middleware.GetUserPhone(r.Context())

    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Protected route accessed",
        "user":    userPhone,
        "session": sessionToken,
    })
}
```

### Using with Fiber Framework

#### Simple Integration
```go
package main

import (
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
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

    // Protected endpoint
    app.Get("/protected", func(c *fiber.Ctx) error {
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

    log.Fatal(app.Listen(":3000"))
}
```

#### Advanced Integration with Middleware
```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
)

func main() {
    app := fiber.New()

    // Protected routes with middleware
    protected := app.Group("/api", fiberAuthMiddleware())
    {
        protected.Get("/protected", func(c *fiber.Ctx) error {
            userPhone := c.Locals("user_phone").(string)
            return c.JSON(fiber.Map{
                "message": "Protected route",
                "user":    userPhone,
            })
        })
    }

    app.Listen(":3000")
}

// Fiber authentication middleware
func fiberAuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Missing authorization header",
            })
        }

        if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid authorization format",
            })
        }

        sessionToken := authHeader[7:]
        userPhone := c.Get("X-User-Phone")
        if userPhone == "" {
            userPhone = c.Query("user_phone")
        }

        if userPhone == "" {
            return c.Status(400).JSON(fiber.Map{
                "error": "Missing user phone",
            })
        }

        // Verify session
        verified, err := rauthprovider.VerifySession(c.Context(), sessionToken, userPhone)
        if err != nil || !verified {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid session",
            })
        }

        // Add to context
        c.Locals("session_token", sessionToken)
        c.Locals("user_phone", userPhone)

        return c.Next()
    }
}
```

### Using with Gin Framework
```go
package main

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
)

func main() {
    // Initialize provider
    config := &rauthprovider.Config{
        RauthAPIKey:   os.Getenv("RAUTH_API_KEY"),
        AppID:         os.Getenv("RAUTH_APP_ID"),
        WebhookSecret: os.Getenv("RAUTH_WEBHOOK_SECRET"),
    }
    rauthprovider.Init(config)

    r := gin.Default()

    // Webhook endpoint
    r.POST("/rauth/webhook", gin.WrapF(rauthprovider.WebhookHandler()))

    // Protected routes
    protected := r.Group("/api")
    protected.Use(ginAuthMiddleware())
    {
        protected.GET("/protected", func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "Protected route"})
        })
    }

    r.Run(":8080")
}

func ginAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        sessionToken := c.GetHeader("Authorization")
        if sessionToken == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization"})
            c.Abort()
            return
        }

        userPhone := c.GetHeader("X-User-Phone")
        if userPhone == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user phone"})
            c.Abort()
            return
        }

        // Remove "Bearer " prefix
        if len(sessionToken) > 7 && sessionToken[:7] == "Bearer " {
            sessionToken = sessionToken[7:]
        }

        verified, err := rauthprovider.VerifySession(c.Request.Context(), sessionToken, userPhone)
        if err != nil || !verified {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
            c.Abort()
            return
        }

        c.Set("session_token", sessionToken)
        c.Set("user_phone", userPhone)
        c.Next()
    }
}
```

---

## API Reference

### Configuration
```go
type Config struct {
    RauthAPIKey      string // Your Rauth API key
    AppID            string // Your Rauth app ID
    WebhookSecret    string // Your webhook secret
    DefaultSessionTTL int    // Session TTL in seconds (default: 900)
    DefaultRevokedTTL int    // Revoked session TTL in seconds (default: 3600)
}
```

### Core Functions

#### `rauthprovider.Init(config *rauthprovider.Config) error`
Initialize the RauthProvider with configuration.

#### `rauthprovider.VerifySession(ctx context.Context, sessionToken, userPhone string) (bool, error)`
Verify if a session is valid and matches the phone number.

#### `rauthprovider.IsSessionRevoked(ctx context.Context, sessionToken string) (bool, error)`
Check if a session has been revoked.

#### `rauthprovider.CheckAPIHealth(ctx context.Context) (bool, error)`
Check if the Rauth API is reachable.

#### `rauthprovider.WebhookHandler() http.HandlerFunc`
Returns HTTP handler for webhook events. Uses simple webhook secret authentication (Node.js compatible).

**Webhook Authentication:**
- **Header**: `x-webhook-secret`
- **Method**: Simple secret comparison (Node.js style)
- **Security**: Webhook secret verification

#### `rauthprovider.GetStats() map[string]interface{}`
Get statistics about the provider.

### Middleware Functions

#### `middleware.AuthMiddleware() func(http.Handler) http.Handler`
Creates middleware that requires valid Rauth authentication.

#### `middleware.OptionalAuthMiddleware() func(http.Handler) http.Handler`
Creates middleware that optionally verifies Rauth authentication.

#### `middleware.GetSessionToken(ctx context.Context) (string, bool)`
Extract session token from request context.

#### `middleware.GetUserPhone(ctx context.Context) (string, bool)`
Extract user phone from request context.

---

## Environment Variables

Create a `.env` file with the following variables:

```env
RAUTH_API_KEY=your-rauth-api-key
RAUTH_APP_ID=your-app-id
RAUTH_WEBHOOK_SECRET=your-webhook-secret
```

---

## Error Handling

The library provides detailed error types:

```go
// Common errors
var (
    ErrNotInitialized     = errors.New("rauth provider not initialized")
    ErrInvalidConfig      = errors.New("invalid configuration")
    ErrSessionNotFound    = errors.New("session not found")
    ErrSessionExpired     = errors.New("session expired")
    ErrSessionRevoked     = errors.New("session revoked")
    ErrInvalidSignature   = errors.New("invalid signature")
    ErrAPIUnreachable     = errors.New("rauth API unreachable")
    ErrInvalidPhoneNumber = errors.New("invalid phone number")
)

// Custom error types
type ConfigError struct {
    Field   string
    Message string
}

type APIError struct {
    StatusCode int
    Message    string
}
```

---

## Architecture

The library follows clean architecture principles:

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
```

### Design Patterns

- **Singleton Pattern**: Thread-safe singleton provider using `sync.Once`
- **Dependency Injection**: Interface-based dependencies for testability
- **Repository Pattern**: Abstract data access layer
- **Middleware Pattern**: Chainable HTTP middleware
- **Factory Pattern**: Component creation and initialization

---

## Testing

```go
package main

import (
    "context"
    "testing"

    "github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
)

func TestSessionVerification(t *testing.T) {
    // Initialize provider
    config := &rauthprovider.Config{
        RauthAPIKey:   "test-key",
        AppID:         "test-app",
        WebhookSecret: "test-secret",
    }
    
    if err := rauthprovider.Init(config); err != nil {
        t.Fatalf("Failed to initialize: %v", err)
    }

    // Test session verification
    verified, err := rauthprovider.VerifySession(context.Background(), "test-token", "+1234567890")
    if err != nil {
        t.Errorf("Verification failed: %v", err)
    }

    t.Logf("Verification result: %v", verified)
}
```

---

## Performance Considerations

- **In-Memory Storage**: Sessions are stored in memory for fast access
- **Automatic Cleanup**: Expired sessions are automatically cleaned up every 5 minutes
- **Thread-Safe**: All operations are thread-safe using read-write mutexes
- **Connection Pooling**: HTTP client uses connection pooling for API calls
- **Context Support**: All operations support context for cancellation and timeouts

---

## Security Features

- **Signature Verification**: Webhook signatures are verified using HMAC-SHA256
- **Session Validation**: Sessions are validated against phone numbers
- **Revocation Tracking**: Revoked sessions are tracked and checked
- **TTL Management**: Automatic expiration of sessions and revoked session records
- **Input Validation**: All inputs are validated before processing

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

---

## License

MIT License - see LICENSE file for details.

---

## Support

For support and questions:
- GitHub Issues: [Create an issue](https://github.com/rauth/rauth-provider/issues)
- Documentation: [Rauth.io Documentation](https://docs.rauth.io)
- Email: support@rauth.io 