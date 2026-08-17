package dto

import (
	"time"

	"ava/internal/model"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Username   string `json:"username" validate:"required,min=3,max=30"`
	Name       string `json:"name" validate:"required"`
	Password   string `json:"password" validate:"required,min=8"`
	Phone      string `json:"phone,omitempty"`
	TenantName string `json:"tenant_name" validate:"required"`
	TenantSlug string `json:"tenant_slug" validate:"required,min=3,max=40"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type UserResponse struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	Username        string     `json:"username"`
	Name            string     `json:"name"`
	Phone           string     `json:"phone,omitempty"`
	EmailVerified   bool       `json:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type AuthResponse struct {
	User   UserResponse   `json:"user"`
	Tokens *TokenResponse `json:"tokens,omitempty"`
}

type TenantSummary struct {
	ID   uuid.UUID        `json:"id"`
	Name string           `json:"name"`
	Slug string           `json:"slug"`
	Role model.TenantRole `json:"role"`
}

type RegisterResponse struct {
	User   UserResponse  `json:"user"`
	Tenant TenantSummary `json:"tenant"`
}

type SessionResponse struct {
	ID         uuid.UUID `json:"id"`
	DeviceName string    `json:"device_name"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	Current    bool      `json:"current"`
}

type DeviceInfo struct {
	DeviceName string
	IPAddress  string
	UserAgent  string
}
