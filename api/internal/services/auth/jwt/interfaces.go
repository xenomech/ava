package jwt

import (
	"time"

	"ava/api/internal/config"
	"ava/api/internal/model"

	"github.com/google/uuid"
)

type TokenManager interface {
	GenerateToken(user *model.User, tenantID uuid.UUID, role model.TenantRole, sessionID uuid.UUID, tokenType TokenType, rid string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	GetAccessExpiry() time.Duration
	GetRefreshExpiry() time.Duration
}

type jwtTokenManager struct {
	secretKey     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewTokenManager() TokenManager {
	cfg := config.GetConfig()

	return &jwtTokenManager{
		secretKey:     cfg.JwtSecretKey,
		accessExpiry:  cfg.JwtAccessExpiry,
		refreshExpiry: cfg.JwtRefreshExpiry,
	}
}
