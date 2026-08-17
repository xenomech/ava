package auth

import (
	"context"

	"ava/internal/dto"
	"ava/internal/model"
	membershiprepo "ava/internal/repository/membership"
	sessionrepo "ava/internal/repository/session"
	tenantrepo "ava/internal/repository/tenant"
	tokenrepo "ava/internal/repository/token"
	userrepo "ava/internal/repository/user"
	"ava/internal/services/auth/jwt"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	VerifyEmail(ctx context.Context, tokenString string) error
	ResendVerification(ctx context.Context, emailAddr string) error
	Login(ctx context.Context, req *dto.LoginRequest, deviceInfo dto.DeviceInfo) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*dto.TokenResponse, error)
	AcceptInvite(ctx context.Context, inviteToken string) (*dto.TenantSummary, error)
	SwitchTenant(ctx context.Context, tenantID, userID, sessionID uuid.UUID, tenantSlug string, deviceInfo dto.DeviceInfo) (*dto.AuthResponse, error)
	Logout(ctx context.Context, tenantID, sessionID uuid.UUID) error
	LogoutAll(ctx context.Context, tenantID, userID uuid.UUID) error
	CurrentSession(ctx context.Context, tenantID, userID uuid.UUID) (*dto.AuthResponse, error)
	ValidateSession(ctx context.Context, tenantID, sessionID uuid.UUID) (*model.Session, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
}

type authService struct {
	userRepo       userrepo.Repository
	tenantRepo     tenantrepo.Repository
	membershipRepo membershiprepo.Repository
	sessionRepo    sessionrepo.Repository
	tokenRepo      tokenrepo.Repository
	tokenManager   jwt.TokenManager
}

func NewService(
	userRepo userrepo.Repository,
	tenantRepo tenantrepo.Repository,
	membershipRepo membershiprepo.Repository,
	sessionRepo sessionrepo.Repository,
	tokenRepo tokenrepo.Repository,
) Service {
	return &authService{
		userRepo:       userRepo,
		tenantRepo:     tenantRepo,
		membershipRepo: membershipRepo,
		sessionRepo:    sessionRepo,
		tokenRepo:      tokenRepo,
		tokenManager:   jwt.NewTokenManager(),
	}
}
