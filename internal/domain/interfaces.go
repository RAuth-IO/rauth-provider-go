package domain

import (
	"context"
)

// SessionRepository defines the interface for session storage operations
type SessionRepository interface {
	// Store stores a session
	Store(ctx context.Context, session *Session) error

	// Get retrieves a session by token
	Get(ctx context.Context, token string) (*Session, error)

	// Delete removes a session
	Delete(ctx context.Context, token string) error

	// Cleanup removes expired sessions
	Cleanup(ctx context.Context) error
}

// RevokedSessionRepository defines the interface for revoked session storage operations
type RevokedSessionRepository interface {
	// Store stores a revoked session
	Store(ctx context.Context, revokedSession *RevokedSession) error

	// Get retrieves a revoked session by token
	Get(ctx context.Context, token string) (*RevokedSession, error)

	// Cleanup removes expired revoked sessions
	Cleanup(ctx context.Context) error
}

// APIClient defines the interface for Rauth API communication
type APIClient interface {
	// VerifySession verifies a session with the Rauth API
	VerifySession(ctx context.Context, sessionToken, userPhone string) (bool, error)

	// CheckHealth checks if the Rauth API is reachable
	CheckHealth(ctx context.Context) (bool, error)
}

// WebhookHandler defines the interface for webhook processing
type WebhookHandler interface {
	// ProcessWebhook processes incoming webhook events
	ProcessWebhook(ctx context.Context, event *WebhookEvent) error

	// VerifySignature verifies the webhook signature
	VerifySignature(ctx context.Context, payload []byte, signature string) (bool, error)
}

// SessionService defines the interface for session business logic
type SessionService interface {
	// VerifySession verifies if a session is valid
	VerifySession(ctx context.Context, sessionToken, userPhone string) (bool, error)

	// IsSessionRevoked checks if a session has been revoked
	IsSessionRevoked(ctx context.Context, sessionToken string) (bool, error)

	// RevokeSession revokes a session
	RevokeSession(ctx context.Context, sessionToken string) error
}
