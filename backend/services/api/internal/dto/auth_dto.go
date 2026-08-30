package dto

import (
	"time"

	"ava/api/internal/model"

	"github.com/google/uuid"
)

// RegisterRequest is deliberately three fields; Register derives the username and home slug itself.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	TenantSlug string `json:"tenant_slug,omitempty"`
}

type SwitchTenantRequest struct {
	TenantSlug string `json:"tenant_slug" validate:"required"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type AcceptInviteRequest struct {
	Token string `json:"token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
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
	Tenant *TenantSummary `json:"tenant,omitempty"`
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
