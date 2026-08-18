package api

import "context"

const (
	CodeAuthorizationPending = "authorization_pending"
	CodeSlowDown             = "slow_down"
	CodeAccessDenied         = "access_denied"
	CodeExpiredToken         = "expired_token"
	CodeInvalidRefreshToken  = "invalid_refresh_token"
	CodeDeviceRevoked        = "device_revoked"
)

type DeviceCode struct {
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

type Device struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type DeviceTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Device       Device `json:"device"`
	Tenant       Tenant `json:"tenant"`
}

func (c *Client) RequestDeviceCode(ctx context.Context, deviceName string) (*DeviceCode, error) {
	var out DeviceCode

	body := map[string]string{"device_name": deviceName}
	if err := c.do(ctx, "POST", "/devices/code", body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) PollDeviceToken(ctx context.Context, deviceCode string) (*DeviceTokens, error) {
	var out DeviceTokens

	body := map[string]string{"device_code": deviceCode}
	if err := c.do(ctx, "POST", "/devices/token", body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) RefreshDeviceToken(ctx context.Context, refreshToken string) (*DeviceTokens, error) {
	var out DeviceTokens

	body := map[string]string{"refresh_token": refreshToken}
	if err := c.do(ctx, "POST", "/devices/token/refresh", body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) Heartbeat(ctx context.Context) error {
	return c.do(ctx, "POST", "/devices/heartbeat", nil, nil)
}

func (c *Client) Health(ctx context.Context) error {
	return c.do(ctx, "GET", "/health", nil, nil)
}
