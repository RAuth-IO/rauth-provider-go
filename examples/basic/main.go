package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RAuth-IO/rauth-provider-go/pkg/middleware"
	"github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
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

	// Create HTTP server
	mux := http.NewServeMux()

	// Webhook endpoint
	mux.HandleFunc("/rauth/webhook", rauthprovider.WebhookHandler())

	// Session verification endpoint
	mux.HandleFunc("/api/login", loginHandler)

	// Protected route with authentication middleware
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/protected", protectedHandler)

	// Apply authentication middleware to protected routes
	authMiddleware := middleware.AuthMiddleware()
	protectedHandler := authMiddleware(protectedMux)

	// Mount protected routes
	mux.Handle("/api/", http.StripPrefix("/api", protectedHandler))

	// Health check endpoint
	mux.HandleFunc("/health", healthHandler)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// loginHandler handles session verification
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionToken string `json:"session_token"`
		UserPhone    string `json:"user_phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify session
	verified, err := rauthprovider.VerifySession(r.Context(), req.SessionToken, req.UserPhone)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "verification_failed", err.Error())
		return
	}

	if !verified {
		middleware.WriteError(w, http.StatusUnauthorized, "invalid_session", "Phone number not verified")
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Session verified successfully",
	})
}

// protectedHandler handles protected routes
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// Get session info from context (set by middleware)
	sessionToken, _ := middleware.GetSessionToken(r.Context())
	userPhone, _ := middleware.GetUserPhone(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Protected route accessed successfully",
		"user":    userPhone,
		"session": sessionToken,
	})
}

// healthHandler handles health checks
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Check API health
	healthy, err := rauthprovider.CheckAPIHealth(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusServiceUnavailable, "health_check_failed", err.Error())
		return
	}

	if !healthy {
		middleware.WriteError(w, http.StatusServiceUnavailable, "api_unhealthy", "Rauth API is not responding")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"api":    "connected",
	})
}
