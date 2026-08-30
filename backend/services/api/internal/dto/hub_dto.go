package dto

import (
	"time"

	"github.com/google/uuid"
)

type DeviceCodeRequest struct {
	HubName string `json:"hub_name" validate:"required,max=100"`
}

// HubHeartbeatRequest reports a hub's own view of itself; a nil BrokerConnected means "not saying".
type HubHeartbeatRequest struct {
	BrokerConnected *bool `json:"broker_connected,omitempty"`
}

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int64  `json:"expires_in"`
	Interval        int64  `json:"interval"`
}

type HubPollRequest struct {
	DeviceCode string `json:"device_code" validate:"required"`
}

type HubRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ActivateHubRequest struct {
	UserCode string `json:"user_code" validate:"required"`
}

type BrokerCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type HubTokenResponse struct {
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	ExpiresIn    int64              `json:"expires_in"`
	Hub          HubResponse        `json:"hub"`
	Tenant       TenantSummary      `json:"tenant"`
	Broker       *BrokerCredentials `json:"broker,omitempty"`
}

type HubResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	Online     bool       `json:"online"`
	LastSeenAt *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}
