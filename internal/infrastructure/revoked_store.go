package infrastructure

import (
	"context"
	"sync"
	"time"

	"github.com/rauth/rauth-provider-go/internal/domain"
)

// RevokedSessionStore implements the domain.RevokedSessionRepository interface
type RevokedSessionStore struct {
	revokedSessions map[string]*domain.RevokedSession
	mutex           sync.RWMutex
}

// NewRevokedSessionStore creates a new revoked session store
func NewRevokedSessionStore() *RevokedSessionStore {
	return &RevokedSessionStore{
		revokedSessions: make(map[string]*domain.RevokedSession),
	}
}

// Store stores a revoked session
func (s *RevokedSessionStore) Store(ctx context.Context, revokedSession *domain.RevokedSession) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.revokedSessions[revokedSession.Token] = revokedSession
	return nil
}

// Get retrieves a revoked session by token
func (s *RevokedSessionStore) Get(ctx context.Context, token string) (*domain.RevokedSession, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	revokedSession, exists := s.revokedSessions[token]
	if !exists {
		return nil, domain.ErrSessionNotFound
	}

	// Check if revoked session record is expired
	if time.Now().After(revokedSession.ExpiresAt) {
		// Remove expired revoked session record
		s.mutex.RUnlock()
		s.mutex.Lock()
		delete(s.revokedSessions, token)
		s.mutex.Unlock()
		s.mutex.RLock()
		return nil, domain.ErrSessionNotFound
	}

	return revokedSession, nil
}

// Cleanup removes expired revoked session records
func (s *RevokedSessionStore) Cleanup(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for token, revokedSession := range s.revokedSessions {
		if now.After(revokedSession.ExpiresAt) {
			delete(s.revokedSessions, token)
		}
	}

	return nil
}

// GetStats returns statistics about the revoked session store
func (s *RevokedSessionStore) GetStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	now := time.Now()
	active := 0
	expired := 0

	for _, revokedSession := range s.revokedSessions {
		if now.After(revokedSession.ExpiresAt) {
			expired++
		} else {
			active++
		}
	}

	return map[string]interface{}{
		"total_revoked_sessions":   len(s.revokedSessions),
		"active_revoked_sessions":  active,
		"expired_revoked_sessions": expired,
	}
}
