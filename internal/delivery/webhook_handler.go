package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RAuth-IO/rauth-provider-go/internal/domain"
)

// WebhookHandler implements the domain.WebhookHandler interface
type WebhookHandler struct {
	webhookSecret  string
	sessionService domain.SessionService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(webhookSecret string, sessionService domain.SessionService) *WebhookHandler {
	return &WebhookHandler{
		webhookSecret:  webhookSecret,
		sessionService: sessionService,
	}
}

// ProcessWebhook processes incoming webhook events
func (h *WebhookHandler) ProcessWebhook(ctx context.Context, event *domain.WebhookEvent) error {
	switch event.Type {
	case "session_verified":
		// Session was verified, no action needed as it's handled during verification
		return nil
	case "session_revoked":
		// Session was revoked, add to revoked sessions
		return h.sessionService.RevokeSession(ctx, event.SessionToken)
	default:
		return fmt.Errorf("unknown webhook event type: %s", event.Type)
	}
}

// HTTPHandler returns an http.HandlerFunc for processing webhook requests
func (h *WebhookHandler) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Get the webhook secret from headers (Node.js style)
		webhookSecret := r.Header.Get("x-webhook-secret")
		if webhookSecret == "" {
			http.Error(w, "Missing webhook secret", http.StatusBadRequest)
			return
		}

		// Verify the webhook secret (Node.js style)
		if webhookSecret != h.webhookSecret {
			http.Error(w, "Invalid webhook secret", http.StatusUnauthorized)
			return
		}

		// Parse the webhook event
		var event domain.WebhookEvent
		if err := json.Unmarshal(body, &event); err != nil {
			http.Error(w, "Failed to parse webhook event", http.StatusBadRequest)
			return
		}

		// Process the webhook event
		if err := h.ProcessWebhook(ctx, &event); err != nil {
			http.Error(w, "Failed to process webhook", http.StatusInternalServerError)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}
