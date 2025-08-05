package domain

import (
	"time"
)

// Session represents a user session in the Rauth system
type Session struct {
	Token     string    `json:"token"`
	UserPhone string    `json:"user_phone"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RevokedSession represents a revoked session
type RevokedSession struct {
	Token     string    `json:"token"`
	RevokedAt time.Time `json:"revoked_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Config holds the configuration for RauthProvider
type Config struct {
	RauthAPIKey       string `json:"rauth_api_key"`
	AppID             string `json:"app_id"`
	WebhookSecret     string `json:"webhook_secret"`
	DefaultSessionTTL int    `json:"default_session_ttl"` // in seconds
	DefaultRevokedTTL int    `json:"default_revoked_ttl"` // in seconds
}

// WebhookEvent represents a webhook event from Rauth.io (Node.js compatible)
type WebhookEvent struct {
	Event        string `json:"event"`        // Node.js uses "event" instead of "type"
	SessionToken string `json:"session_token"`
	Phone        string `json:"phone"`        // Node.js uses "phone" instead of "user_phone"
	TTL          int    `json:"ttl"`          // Node.js uses "ttl" instead of "timestamp"
	Reason       string `json:"reason"`       // Node.js specific field
	Signature    string `json:"signature"`
	
	// Legacy fields for backward compatibility
	Type         string `json:"type,omitempty"`
	UserPhone    string `json:"user_phone,omitempty"`
	Timestamp    int64  `json:"timestamp,omitempty"`
}

// APIResponse represents a response from the Rauth API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
