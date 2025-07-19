package rauthprovider

// Config holds the configuration for RauthProvider
type Config struct {
	RauthAPIKey       string `json:"rauth_api_key"`
	AppID             string `json:"app_id"`
	WebhookSecret     string `json:"webhook_secret"`
	DefaultSessionTTL int    `json:"default_session_ttl,omitempty"`
	DefaultRevokedTTL int    `json:"default_revoked_ttl,omitempty"`
}
