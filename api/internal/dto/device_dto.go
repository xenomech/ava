package dto

import (
	"time"

	"github.com/google/uuid"
)

type DeviceCodeRequest struct {
	DeviceName string `json:"device_name" validate:"required,max=100"`
}

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int64  `json:"expires_in"`
	Interval        int64  `json:"interval"`
}

type DevicePollRequest struct {
	DeviceCode string `json:"device_code" validate:"required"`
}

type DeviceRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ActivateDeviceRequest struct {
	UserCode string `json:"user_code" validate:"required"`
}

type DeviceTokenResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    int64          `json:"expires_in"`
	Device       DeviceResponse `json:"device"`
	Tenant       TenantSummary  `json:"tenant"`
}

type DeviceResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	LastSeenAt *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}
