package dto

import (
	"time"

	"ava/api/internal/model"

	"github.com/google/uuid"
)

type CreateAPITokenRequest struct {
	Name   string   `json:"name" validate:"required,max=80"`
	Scopes []string `json:"scopes" validate:"required,min=1,max=16,dive,required"`
	// ExpiresInDays is optional; omit it for a token that does not expire.
	ExpiresInDays int `json:"expires_in_days" validate:"omitempty,min=1,max=3650"`
}

type APITokenResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Scopes     []string   `json:"scopes"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// CreatedAPITokenResponse carries the plaintext, which is shown once and never retrievable again.
type CreatedAPITokenResponse struct {
	Token APITokenResponse `json:"token"`
	Value string           `json:"value"`
}

func NewAPITokenResponse(token *model.APIToken) APITokenResponse {
	return APITokenResponse{
		ID:         token.ID,
		Name:       token.Name,
		Scopes:     []string(token.Scopes),
		LastUsedAt: token.LastUsedAt,
		ExpiresAt:  token.ExpiresAt,
		RevokedAt:  token.RevokedAt,
		CreatedAt:  token.CreatedAt,
	}
}

type ScopeResponse struct {
	Scopes []string `json:"scopes"`
}
