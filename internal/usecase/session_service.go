package usecase

import (
	"context"
	"time"

	"github.com/RAuth-IO/rauth-provider-go/internal/domain"
)

// SessionService implements the domain.SessionService interface
type SessionService struct {
	sessionRepo        domain.SessionRepository
	revokedSessionRepo domain.RevokedSessionRepository
	apiClient          domain.APIClient
	config             *domain.Config
}

// NewSessionService creates a new session service
func NewSessionService(
	sessionRepo domain.SessionRepository,
	revokedSessionRepo domain.RevokedSessionRepository,
	apiClient domain.APIClient,
	config *domain.Config,
) *SessionService {
	return &SessionService{
		sessionRepo:        sessionRepo,
		revokedSessionRepo: revokedSessionRepo,
		apiClient:          apiClient,
		config:             config,
	}
}

// VerifySession verifies if a session is valid
func (s *SessionService) VerifySession(ctx context.Context, sessionToken, userPhone string) (bool, error) {
	// First check if session is revoked
	isRevoked, err := s.IsSessionRevoked(ctx, sessionToken)
	if err != nil && err != domain.ErrSessionNotFound {
		return false, err
	}
	if isRevoked {
		return false, domain.ErrSessionRevoked
	}

	// Check local session store first
	session, err := s.sessionRepo.Get(ctx, sessionToken)
	if err == nil {
		// Session found locally, verify phone number matches
		if session.UserPhone == userPhone {
			return true, nil
		}
		return false, domain.ErrInvalidPhoneNumber
	}

	// Session not found locally, check with API
	verified, err := s.apiClient.VerifySession(ctx, sessionToken, userPhone)
	if err != nil {
		return false, err
	}

	if verified {
		// Store the verified session locally
		now := time.Now()
		session := &domain.Session{
			Token:     sessionToken,
			UserPhone: userPhone,
			CreatedAt: now,
			ExpiresAt: now.Add(time.Duration(s.config.DefaultSessionTTL) * time.Second),
		}

		if err := s.sessionRepo.Store(ctx, session); err != nil {
			// Log error but don't fail the verification
			// The session is still valid according to the API
		}
	}

	return verified, nil
}

// IsSessionRevoked checks if a session has been revoked
func (s *SessionService) IsSessionRevoked(ctx context.Context, sessionToken string) (bool, error) {
	_, err := s.revokedSessionRepo.Get(ctx, sessionToken)
	if err == domain.ErrSessionNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// RevokeSession revokes a session
func (s *SessionService) RevokeSession(ctx context.Context, sessionToken string) error {
	// Remove from active sessions
	if err := s.sessionRepo.Delete(ctx, sessionToken); err != nil && err != domain.ErrSessionNotFound {
		return err
	}

	// Add to revoked sessions
	now := time.Now()
	revokedSession := &domain.RevokedSession{
		Token:     sessionToken,
		RevokedAt: now,
		ExpiresAt: now.Add(time.Duration(s.config.DefaultRevokedTTL) * time.Second),
	}

	return s.revokedSessionRepo.Store(ctx, revokedSession)
}

// Cleanup performs cleanup of expired sessions and revoked sessions
func (s *SessionService) Cleanup(ctx context.Context) error {
	// Cleanup expired sessions
	if err := s.sessionRepo.Cleanup(ctx); err != nil {
		return err
	}

	// Cleanup expired revoked sessions
	return s.revokedSessionRepo.Cleanup(ctx)
}
