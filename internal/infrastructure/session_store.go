package infrastructure

import (
	"context"
	"sync"
	"time"

	"github.com/rauth/rauth-provider-go/internal/domain"
)

// SessionStore implements the domain.SessionRepository interface
type SessionStore struct {
	sessions map[string]*domain.Session
	mutex    sync.RWMutex
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*domain.Session),
	}
}

// Store stores a session
func (s *SessionStore) Store(ctx context.Context, session *domain.Session) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.sessions[session.Token] = session
	return nil
}

// Get retrieves a session by token
func (s *SessionStore) Get(ctx context.Context, token string) (*domain.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, exists := s.sessions[token]
	if !exists {
		return nil, domain.ErrSessionNotFound
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Remove expired session
		s.mutex.RUnlock()
		s.mutex.Lock()
		delete(s.sessions, token)
		s.mutex.Unlock()
		s.mutex.RLock()
		return nil, domain.ErrSessionExpired
	}

	return session, nil
}

// Delete removes a session
func (s *SessionStore) Delete(ctx context.Context, token string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.sessions, token)
	return nil
}

// Cleanup removes expired sessions
func (s *SessionStore) Cleanup(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for token, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, token)
		}
	}

	return nil
}

// GetStats returns statistics about the session store
func (s *SessionStore) GetStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	now := time.Now()
	active := 0
	expired := 0

	for _, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			expired++
		} else {
			active++
		}
	}

	return map[string]interface{}{
		"total_sessions":   len(s.sessions),
		"active_sessions":  active,
		"expired_sessions": expired,
	}
}
