// Package rauthprovider provides a lightweight, plug-and-play Go library for phone number authentication
// using the Rauth.io reverse verification flow via WhatsApp or SMS.
//
// Features:
// - Reverse Authentication via WhatsApp or SMS
// - Session Management with TTL
// - Webhook Support for real-time updates
// - HTTP Middleware integration
// - Signature-based verification
// - In-memory session tracking with API fallback
// - Singleton pattern for state management
package rauthprovider

import (
	"context"
	"net/http"

	"github.com/rauth/rauth-provider-go/internal/domain"
)

// Provider is the main interface for Rauth authentication
type Provider interface {
	// Init initializes the provider with configuration
	Init(config *domain.Config) error

	// VerifySession verifies if a session is valid
	VerifySession(ctx context.Context, sessionToken, userPhone string) (bool, error)

	// IsSessionRevoked checks if a session has been revoked
	IsSessionRevoked(ctx context.Context, sessionToken string) (bool, error)

	// CheckAPIHealth checks if the Rauth API is reachable
	CheckAPIHealth(ctx context.Context) (bool, error)

	// WebhookHandler returns the HTTP handler for webhook processing
	WebhookHandler() http.HandlerFunc

	// GetStats returns statistics about the provider
	GetStats() map[string]interface{}
}

// GetProvider returns the singleton provider instance
func GetProvider() Provider {
	return GetInstance()
}

// Init is a convenience function to initialize the provider
func Init(config *domain.Config) error {
	return GetInstance().Init(config)
}

// VerifySession is a convenience function to verify a session
func VerifySession(ctx context.Context, sessionToken, userPhone string) (bool, error) {
	return GetInstance().VerifySession(ctx, sessionToken, userPhone)
}

// IsSessionRevoked is a convenience function to check if a session is revoked
func IsSessionRevoked(ctx context.Context, sessionToken string) (bool, error) {
	return GetInstance().IsSessionRevoked(ctx, sessionToken)
}

// CheckAPIHealth is a convenience function to check API health
func CheckAPIHealth(ctx context.Context) (bool, error) {
	return GetInstance().CheckAPIHealth(ctx)
}

// WebhookHandler is a convenience function to get the webhook handler
func WebhookHandler() http.HandlerFunc {
	return GetInstance().WebhookHandler()
}

// GetStats is a convenience function to get provider statistics
func GetStats() map[string]interface{} {
	return GetInstance().GetStats()
}
