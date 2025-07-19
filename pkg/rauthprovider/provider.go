package rauthprovider

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/RAuth-IO/rauth-provider-go/internal/delivery"
	"github.com/RAuth-IO/rauth-provider-go/internal/domain"
	"github.com/RAuth-IO/rauth-provider-go/internal/infrastructure"
	"github.com/RAuth-IO/rauth-provider-go/internal/usecase"
)

// RauthProvider is the main provider for Rauth authentication
type RauthProvider struct {
	config         *Config
	sessionService *usecase.SessionService
	apiClient      *infrastructure.APIClient
	webhookHandler *delivery.WebhookHandler
	initialized    bool
	mutex          sync.RWMutex
}

var (
	instance *RauthProvider
	once     sync.Once
)

// GetInstance returns the singleton instance of RauthProvider
func GetInstance() *RauthProvider {
	once.Do(func() {
		instance = &RauthProvider{}
	})
	return instance
}

// Init initializes the RauthProvider with configuration
func (p *RauthProvider) Init(config *Config) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Validate configuration
	if err := p.validateConfig(config); err != nil {
		return err
	}

	// Set default values if not provided
	if config.DefaultSessionTTL == 0 {
		config.DefaultSessionTTL = 900 // 15 minutes
	}
	if config.DefaultRevokedTTL == 0 {
		config.DefaultRevokedTTL = 3600 // 1 hour
	}

	// Create infrastructure components
	sessionStore := infrastructure.NewSessionStore()
	revokedSessionStore := infrastructure.NewRevokedSessionStore()
	apiClient := infrastructure.NewAPIClient(config.RauthAPIKey, config.AppID)

	// Convert public config to domain config
	domainConfig := &domain.Config{
		RauthAPIKey:      config.RauthAPIKey,
		AppID:            config.AppID,
		WebhookSecret:    config.WebhookSecret,
		DefaultSessionTTL: config.DefaultSessionTTL,
		DefaultRevokedTTL: config.DefaultRevokedTTL,
	}

	// Create use case layer
	sessionService := usecase.NewSessionService(sessionStore, revokedSessionStore, apiClient, domainConfig)

	// Create webhook handler
	webhookHandler := delivery.NewWebhookHandler(config.WebhookSecret, sessionService)

	// Set the components
	p.config = config
	p.sessionService = sessionService
	p.apiClient = apiClient
	p.webhookHandler = webhookHandler
	p.initialized = true

	// Start cleanup goroutine
	go p.startCleanupRoutine()

	return nil
}

// validateConfig validates the configuration
func (p *RauthProvider) validateConfig(config *Config) error {
	if config == nil {
		return &domain.ConfigError{Field: "config", Message: "configuration cannot be nil"}
	}
	if config.RauthAPIKey == "" {
		return &domain.ConfigError{Field: "rauth_api_key", Message: "rauth API key is required"}
	}
	if config.AppID == "" {
		return &domain.ConfigError{Field: "app_id", Message: "app ID is required"}
	}
	if config.WebhookSecret == "" {
		return &domain.ConfigError{Field: "webhook_secret", Message: "webhook secret is required"}
	}
	return nil
}

// VerifySession verifies if a session is valid
func (p *RauthProvider) VerifySession(ctx context.Context, sessionToken, userPhone string) (bool, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if !p.initialized {
		return false, domain.ErrNotInitialized
	}

	return p.sessionService.VerifySession(ctx, sessionToken, userPhone)
}

// IsSessionRevoked checks if a session has been revoked
func (p *RauthProvider) IsSessionRevoked(ctx context.Context, sessionToken string) (bool, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if !p.initialized {
		return false, domain.ErrNotInitialized
	}

	return p.sessionService.IsSessionRevoked(ctx, sessionToken)
}

// CheckAPIHealth checks if the Rauth API is reachable
func (p *RauthProvider) CheckAPIHealth(ctx context.Context) (bool, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if !p.initialized {
		return false, domain.ErrNotInitialized
	}

	return p.apiClient.CheckHealth(ctx)
}

// WebhookHandler returns the HTTP handler for webhook processing
func (p *RauthProvider) WebhookHandler() http.HandlerFunc {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if !p.initialized {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Provider not initialized", http.StatusInternalServerError)
		}
	}

	return p.webhookHandler.HTTPHandler()
}

// startCleanupRoutine starts a background routine to clean up expired sessions
func (p *RauthProvider) startCleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute) // Run cleanup every 5 minutes
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		if err := p.sessionService.Cleanup(ctx); err != nil {
			// Log error but continue
			// In a production environment, you might want to use a proper logger
		}
	}
}

// GetStats returns statistics about the provider
func (p *RauthProvider) GetStats() map[string]interface{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if !p.initialized {
		return map[string]interface{}{
			"initialized": false,
		}
	}

	return map[string]interface{}{
		"initialized": true,
		"config": map[string]interface{}{
			"app_id":              p.config.AppID,
			"default_session_ttl": p.config.DefaultSessionTTL,
			"default_revoked_ttl": p.config.DefaultRevokedTTL,
		},
	}
}
