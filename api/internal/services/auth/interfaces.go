package auth

import (
	"context"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	membershiprepo "ava/api/internal/repository/membership"
	sessionrepo "ava/api/internal/repository/session"
	tenantrepo "ava/api/internal/repository/tenant"
	tokenrepo "ava/api/internal/repository/token"
	userrepo "ava/api/internal/repository/user"
	"ava/api/internal/services/auth/jwt"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	VerifyEmail(ctx context.Context, tokenString string) error
	ResendVerification(ctx context.Context, emailAddr string) error
	Login(ctx context.Context, req *dto.LoginRequest, deviceInfo dto.DeviceInfo) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*dto.TokenResponse, error)
	ForgotPassword(ctx context.Context, emailAddr string) error
	ResetPassword(ctx context.Context, tokenString, newPassword string) error
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	AcceptInvite(ctx context.Context, inviteToken string) (*dto.TenantSummary, error)
	SwitchTenant(ctx context.Context, tenantID, userID, sessionID uuid.UUID, tenantSlug string, deviceInfo dto.DeviceInfo) (*dto.AuthResponse, error)
	Logout(ctx context.Context, tenantID, sessionID uuid.UUID) error
	LogoutAll(ctx context.Context, tenantID, userID uuid.UUID) error
	GetUserSessions(ctx context.Context, tenantID, userID, currentSessionID uuid.UUID) ([]*dto.SessionResponse, error)
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
