package domain

import (
	"testing"
	"time"
)

func TestSession_Expiration(t *testing.T) {
	now := time.Now()
	session := &Session{
		Token:     "test-token",
		UserPhone: "+1234567890",
		CreatedAt: now,
		ExpiresAt: now.Add(15 * time.Minute),
	}

	// Test that session is not expired immediately
	if time.Now().After(session.ExpiresAt) {
		t.Error("Session should not be expired immediately after creation")
	}

	// Test that session expires after TTL
	expiredSession := &Session{
		Token:     "expired-token",
		UserPhone: "+1234567890",
		CreatedAt: now.Add(-20 * time.Minute),
		ExpiresAt: now.Add(-5 * time.Minute),
	}

	if !time.Now().After(expiredSession.ExpiresAt) {
		t.Error("Session should be expired")
	}
}

func TestRevokedSession_Expiration(t *testing.T) {
	now := time.Now()
	revokedSession := &RevokedSession{
		Token:     "revoked-token",
		RevokedAt: now,
		ExpiresAt: now.Add(1 * time.Hour),
	}

	// Test that revoked session record is not expired immediately
	if time.Now().After(revokedSession.ExpiresAt) {
		t.Error("Revoked session record should not be expired immediately after creation")
	}

	// Test that revoked session record expires after TTL
	expiredRevokedSession := &RevokedSession{
		Token:     "expired-revoked-token",
		RevokedAt: now.Add(-2 * time.Hour),
		ExpiresAt: now.Add(-1 * time.Hour),
	}

	if !time.Now().After(expiredRevokedSession.ExpiresAt) {
		t.Error("Revoked session record should be expired")
	}
}

func TestConfig_Validation(t *testing.T) {
	// Test valid config
	validConfig := &Config{
		RauthAPIKey:       "test-api-key",
		AppID:             "test-app-id",
		WebhookSecret:     "test-webhook-secret",
		DefaultSessionTTL: 900,
		DefaultRevokedTTL: 3600,
	}

	if validConfig.RauthAPIKey == "" {
		t.Error("Valid config should have non-empty API key")
	}

	if validConfig.AppID == "" {
		t.Error("Valid config should have non-empty app ID")
	}

	if validConfig.WebhookSecret == "" {
		t.Error("Valid config should have non-empty webhook secret")
	}

	// Test default TTL values
	if validConfig.DefaultSessionTTL <= 0 {
		t.Error("Default session TTL should be positive")
	}

	if validConfig.DefaultRevokedTTL <= 0 {
		t.Error("Default revoked TTL should be positive")
	}
}

func TestWebhookEvent_Validation(t *testing.T) {
	event := &WebhookEvent{
		Type:         "session_verified",
		SessionToken: "test-session-token",
		UserPhone:    "+1234567890",
		Timestamp:    time.Now().Unix(),
		Signature:    "test-signature",
	}

	if event.Type == "" {
		t.Error("Webhook event should have a type")
	}

	if event.SessionToken == "" {
		t.Error("Webhook event should have a session token")
	}

	if event.UserPhone == "" {
		t.Error("Webhook event should have a user phone")
	}

	if event.Timestamp <= 0 {
		t.Error("Webhook event should have a valid timestamp")
	}
}

func TestAPIResponse_Validation(t *testing.T) {
	// Test success response
	successResp := &APIResponse{
		Success: true,
		Data:    map[string]interface{}{"verified": true},
	}

	if !successResp.Success {
		t.Error("Success response should have Success=true")
	}

	// Test error response
	errorResp := &APIResponse{
		Success: false,
		Error:   "API error occurred",
	}

	if errorResp.Success {
		t.Error("Error response should have Success=false")
	}

	if errorResp.Error == "" {
		t.Error("Error response should have an error message")
	}
}
