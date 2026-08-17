package auth

import (
	"context"

	"ava/internal/dto"
	"ava/internal/model"
	membershiprepo "ava/internal/repository/membership"
	sessionrepo "ava/internal/repository/session"
	tenantrepo "ava/internal/repository/tenant"
	userrepo "ava/internal/repository/user"
	"ava/internal/services/auth/jwt"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest, deviceInfo dto.DeviceInfo) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*dto.TokenResponse, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	CurrentSession(ctx context.Context, userID uuid.UUID) (*dto.AuthResponse, error)
	ValidateSession(ctx context.Context, sessionID uuid.UUID) (*model.Session, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
}

type authService struct {
	userRepo       userrepo.Repository
	tenantRepo     tenantrepo.Repository
	membershipRepo membershiprepo.Repository
	sessionRepo    sessionrepo.Repository
	tokenManager   jwt.TokenManager
}

func NewService(
	userRepo userrepo.Repository,
	tenantRepo tenantrepo.Repository,
	membershipRepo membershiprepo.Repository,
	sessionRepo sessionrepo.Repository,
) Service {
	return &authService{
		userRepo:       userRepo,
		tenantRepo:     tenantRepo,
		membershipRepo: membershipRepo,
		sessionRepo:    sessionRepo,
		tokenManager:   jwt.NewTokenManager(),
	}
}
