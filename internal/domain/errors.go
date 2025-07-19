package domain

import (
	"errors"
	"fmt"
)

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

// ConfigError represents configuration-related errors
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error in field '%s': %s", e.Field, e.Message)
}

func (e *ConfigError) Unwrap() error {
	return ErrInvalidConfig
}

// APIError represents API-related errors
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

func (e *APIError) Unwrap() error {
	return ErrAPIUnreachable
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}
