package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RAuth-IO/rauth-provider-go/internal/domain"
)

// APIClient implements the domain.APIClient interface
type APIClient struct {
	baseURL    string
	apiKey     string
	appID      string
	httpClient *http.Client
}

// NewAPIClient creates a new API client
func NewAPIClient(apiKey, appID string) *APIClient {
	return &APIClient{
		baseURL: "https://api.rauth.io/session",
		apiKey:  apiKey,
		appID:   appID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// VerifySession verifies a session with the Rauth API
func (c *APIClient) VerifySession(ctx context.Context, sessionToken, userPhone string) (bool, error) {
	payload := map[string]interface{}{
		"session_token": sessionToken,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Retry with exponential backoff for Cloudflare challenges
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry with exponential backoff
			backoff := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-time.After(backoff):
			}
		}

		req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/status", bytes.NewBuffer(jsonData))
		if err != nil {
			return false, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("X-App-ID", c.appID)
		req.Header.Set("User-Agent", "RauthProvider-Go/1.0")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return false, &domain.APIError{
				StatusCode: 0,
				Message:    fmt.Sprintf("failed to make request: %v", err),
			}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode == 404 {
			return false, nil // Session not found
		}

		// If we get a 403, retry (Cloudflare challenge)
		if resp.StatusCode == 403 && attempt < maxRetries-1 {
			continue
		}

		if resp.StatusCode == 403 {
			return false, &domain.APIError{
				StatusCode: resp.StatusCode,
				Message:    "Access denied. This might be due to Cloudflare protection. Please check your API key and app ID.",
			}
		}

		if resp.StatusCode != http.StatusOK {
			return false, &domain.APIError{
				StatusCode: resp.StatusCode,
				Message:    string(body),
			}
		}

		var sessionDetails map[string]interface{}
		if err := json.Unmarshal(body, &sessionDetails); err != nil {
			return false, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		// Check if session is verified
		if status, exists := sessionDetails["status"]; exists {
			if statusStr, ok := status.(string); ok {
				if statusStr == "verified" {
					// Verify phone number matches if provided
					if userPhone != "" {
						if phone, exists := sessionDetails["phone"]; exists {
							if phoneStr, ok := phone.(string); ok {
								if phoneStr != userPhone {
									return false, nil
								}
							}
						}
					}
					return true, nil
				}
			}
		}

		return false, nil
	}

	// If we get here, all retries failed
	return false, &domain.APIError{
		StatusCode: 0,
		Message:    "All retry attempts failed",
	}
}

// CheckHealth checks if the Rauth API is reachable
func (c *APIClient) CheckHealth(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("X-App-ID", c.appID)
	req.Header.Set("User-Agent", "RauthProvider-Go/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, &domain.APIError{
			StatusCode: 0,
			Message:    fmt.Sprintf("health check failed: %v", err),
		}
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
