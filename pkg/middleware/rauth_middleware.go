// Package middleware provides HTTP middleware for Rauth authentication
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RAuth-IO/rauth-provider-go/pkg/rauthprovider"
)

// AuthMiddleware creates middleware that verifies Rauth sessions
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract session token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			// Check if it's a Bearer token
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			sessionToken := strings.TrimPrefix(authHeader, "Bearer ")

			// Extract user phone from header or query parameter
			userPhone := r.Header.Get("X-User-Phone")
			if userPhone == "" {
				userPhone = r.URL.Query().Get("user_phone")
			}

			if userPhone == "" {
				http.Error(w, "Missing user phone", http.StatusBadRequest)
				return
			}

			// Verify session
			verified, err := rauthprovider.VerifySession(r.Context(), sessionToken, userPhone)
			if err != nil {
				http.Error(w, "Session verification failed", http.StatusUnauthorized)
				return
			}

			if !verified {
				http.Error(w, "Invalid session", http.StatusUnauthorized)
				return
			}

			// Check if session is revoked
			isRevoked, err := rauthprovider.IsSessionRevoked(r.Context(), sessionToken)
			if err != nil {
				http.Error(w, "Session status check failed", http.StatusInternalServerError)
				return
			}

			if isRevoked {
				http.Error(w, "Session revoked", http.StatusUnauthorized)
				return
			}

			// Add session info to context
			ctx := context.WithValue(r.Context(), "session_token", sessionToken)
			ctx = context.WithValue(ctx, "user_phone", userPhone)

			// Call next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware creates middleware that optionally verifies Rauth sessions
// If authentication fails, it continues without session info in context
func OptionalAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract session token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No auth header, continue without session
				next.ServeHTTP(w, r)
				return
			}

			// Check if it's a Bearer token
			if !strings.HasPrefix(authHeader, "Bearer ") {
				// Invalid format, continue without session
				next.ServeHTTP(w, r)
				return
			}

			sessionToken := strings.TrimPrefix(authHeader, "Bearer ")

			// Extract user phone from header or query parameter
			userPhone := r.Header.Get("X-User-Phone")
			if userPhone == "" {
				userPhone = r.URL.Query().Get("user_phone")
			}

			if userPhone == "" {
				// No user phone, continue without session
				next.ServeHTTP(w, r)
				return
			}

			// Try to verify session
			verified, err := rauthprovider.VerifySession(r.Context(), sessionToken, userPhone)
			if err != nil || !verified {
				// Verification failed, continue without session
				next.ServeHTTP(w, r)
				return
			}

			// Check if session is revoked
			isRevoked, err := rauthprovider.IsSessionRevoked(r.Context(), sessionToken)
			if err != nil || isRevoked {
				// Session revoked or error, continue without session
				next.ServeHTTP(w, r)
				return
			}

			// Add session info to context
			ctx := context.WithValue(r.Context(), "session_token", sessionToken)
			ctx = context.WithValue(ctx, "user_phone", userPhone)

			// Call next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetSessionToken extracts session token from context
func GetSessionToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value("session_token").(string)
	return token, ok
}

// GetUserPhone extracts user phone from context
func GetUserPhone(ctx context.Context) (string, bool) {
	phone, ok := ctx.Value("user_phone").(string)
	return phone, ok
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// WriteError writes an error response as JSON
func WriteError(w http.ResponseWriter, statusCode int, error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:   error,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
