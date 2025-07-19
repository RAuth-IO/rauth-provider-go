# Fiber + RauthProvider Integration Examples

This directory contains examples of how to integrate the RauthProvider Go library with the Fiber web framework.

## 📁 Directory Structure

```
examples/fiber/
├── simple/
│   ├── main.go          # Simple integration example
│   └── go.mod           # Go module for simple example
├── advanced/
│   ├── main.go          # Advanced integration with middleware
│   └── go.mod           # Go module for advanced example
├── README.md            # This file
├── test_fiber.sh        # Bash test script (Linux/Mac)
└── test_fiber.ps1       # PowerShell test script (Windows)
```

## 🚀 Quick Start

### 1. Simple Example

```bash
cd simple
go mod tidy
go run .
```

### 2. Advanced Example

```bash
cd advanced
go mod tidy
go run .
```

### 3. Set Environment Variables

```bash
export RAUTH_API_KEY="your_api_key_here"
export RAUTH_APP_ID="your_app_id_here"
export RAUTH_WEBHOOK_SECRET="your_webhook_secret_here"
```

## 🔧 Features Demonstrated

### Simple Example (`simple/main.go`)
- ✅ Basic session verification
- ✅ Protected routes
- ✅ Health check endpoint
- ✅ Error handling

### Advanced Example (`advanced/main.go`)
- ✅ Custom middleware for authentication
- ✅ Optional authentication middleware
- ✅ Webhook handling
- ✅ Session management (login/logout)
- ✅ Route grouping
- ✅ CORS configuration
- ✅ Custom error handling
- ✅ Comprehensive logging

## 🧪 Testing

### Using Test Scripts

**Linux/Mac:**
```bash
chmod +x test_fiber.sh
./test_fiber.sh
```

**Windows:**
```powershell
.\test_fiber.ps1
```

### Manual Testing

1. **Health Check:**
```bash
curl http://localhost:3000/health
```

2. **Login (will fail with test data):**
```bash
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{"session_token": "test_token", "user_phone": "+1234567890"}'
```

3. **Protected Route (will fail without auth):**
```bash
curl http://localhost:3000/protected
```

4. **Protected Route with Auth (will fail with invalid token):**
```bash
curl -H "Authorization: Bearer invalid_token" \
  -H "X-User-Phone: +1234567890" \
  http://localhost:3000/protected
```

## 🔐 Authentication Flow

1. **Client sends session token and phone number**
2. **Server verifies session using RauthProvider**
3. **If valid, user is authenticated**
4. **Session info is stored in Fiber context**
5. **Protected routes can access user data**

## 📡 Webhook Integration

The advanced example includes webhook handling for real-time session updates:

```go
// Webhook endpoint
app.Post("/rauth/webhook", rauthWebhookHandler)
```

This allows your application to receive real-time updates about session changes from the Rauth service.

## 🔄 Middleware Types

### Required Authentication
```go
protected := app.Group("/api", fiberAuthMiddleware())
```

### Optional Authentication
```go
optionalAuth := app.Group("/api/optional", fiberOptionalAuthMiddleware())
```

## 🛠️ Customization

### Custom Error Responses
```go
func customErrorHandler(c *fiber.Ctx, err error) error {
    return c.Status(500).JSON(fiber.Map{
        "error": true,
        "message": err.Error(),
    })
}
```

### Custom Middleware
```go
func myCustomMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Your custom logic here
        return c.Next()
    }
}
```

## 📊 API Endpoints

### Simple Example
- `GET /health` - Health check
- `POST /login` - Session verification
- `GET /protected` - Protected route

### Advanced Example
- `GET /health` - Health check
- `POST /rauth/webhook` - Webhook endpoint
- `POST /api/login` - Session verification
- `GET /api/protected` - Protected route
- `GET /api/user/profile` - User profile
- `POST /api/logout` - Logout
- `GET /api/optional/public` - Public route (optional auth)
- `GET /api/optional/user-info` - User info (optional auth)

## 🔧 Configuration

The examples use environment variables for configuration:

```go
config := &rauthprovider.Config{
    RauthAPIKey:      os.Getenv("RAUTH_API_KEY"),
    AppID:            os.Getenv("RAUTH_APP_ID"),
    WebhookSecret:    os.Getenv("RAUTH_WEBHOOK_SECRET"),
    DefaultSessionTTL: 900,  // 15 minutes
    DefaultRevokedTTL: 3600, // 1 hour
}
```

## 🚨 Error Handling

The examples include comprehensive error handling:

- Invalid session tokens
- Missing authentication headers
- API communication errors
- Webhook signature verification
- Session revocation checks

## 📈 Performance

- Uses Fiber's high-performance HTTP framework
- In-memory session caching
- Efficient middleware chain
- Minimal overhead for authentication checks

## 🔒 Security

- HMAC-SHA256 signature verification for webhooks
- Session token validation
- Phone number verification
- Session revocation support
- Secure header handling
