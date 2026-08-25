package api

import "context"

const (
	CodeAuthorizationPending = "authorization_pending"
	CodeSlowDown             = "slow_down"
	CodeAccessDenied         = "access_denied"
	CodeExpiredToken         = "expired_token"
	CodeInvalidRefreshToken  = "invalid_refresh_token"
	CodeHubRevoked           = "hub_revoked"
)

type ActivationCode struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int64  `json:"expires_in"`
	Interval        int64  `json:"interval"`
}

type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Hub struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type BrokerCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type HubTokens struct {
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	ExpiresIn    int64              `json:"expires_in"`
	Hub          Hub                `json:"hub"`
	Tenant       Tenant             `json:"tenant"`
	Broker       *BrokerCredentials `json:"broker,omitempty"`
}

func (c *Client) RequestActivationCode(ctx context.Context, hubName string) (*ActivationCode, error) {
	var out ActivationCode

	body := map[string]string{"hub_name": hubName}
	if err := c.do(ctx, "POST", "/hubs/code", body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) PollToken(ctx context.Context, deviceCode string) (*HubTokens, error) {
	var out HubTokens

	body := map[string]string{"device_code": deviceCode}
	if err := c.do(ctx, "POST", "/hubs/token", body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*HubTokens, error) {
	var out HubTokens

	body := map[string]string{"refresh_token": refreshToken}
	if err := c.do(ctx, "POST", "/hubs/token/refresh", body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// Heartbeat reports that the hub is alive, and whether commands can still
// reach it. The second half is the useful one: HTTP staying up says nothing
// about the broker, and for an hour it said exactly that while every command
// went nowhere.
func (c *Client) Heartbeat(ctx context.Context, brokerConnected bool) error {
	body := map[string]bool{"broker_connected": brokerConnected}

	return c.do(ctx, "POST", "/hubs/heartbeat", body, nil)
}

func (c *Client) Health(ctx context.Context) error {
	return c.do(ctx, "GET", "/health", nil, nil)
}
