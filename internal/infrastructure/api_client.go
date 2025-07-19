package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rauth/rauth-provider-go/internal/domain"
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
		baseURL: "https://api.rauth.io",
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
		"user_phone":    userPhone,
		"app_id":        c.appID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/verify-session", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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

	if resp.StatusCode != http.StatusOK {
		return false, &domain.APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	var apiResp domain.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return false, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !apiResp.Success {
		return false, fmt.Errorf("API error: %s", apiResp.Error)
	}

	// Check if the data indicates verification success
	if data, ok := apiResp.Data.(map[string]interface{}); ok {
		if verified, exists := data["verified"]; exists {
			if verifiedBool, ok := verified.(bool); ok {
				return verifiedBool, nil
			}
		}
	}

	return false, fmt.Errorf("unexpected response format")
}

// CheckHealth checks if the Rauth API is reachable
func (c *APIClient) CheckHealth(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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
