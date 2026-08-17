package auth

import (
	"context"
	"time"

	"ava/internal/dto"
	"ava/internal/model"
	sessionrepo "ava/internal/repository/session"
	tenantrepo "ava/internal/repository/tenant"
	userrepo "ava/internal/repository/user"
	"ava/internal/services/auth/jwt"
	"ava/pkg/logger"
	"ava/pkg/serrors"

	"github.com/google/uuid"
)

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	if !model.IsValidSlug(req.TenantSlug) {
		return nil, ErrInvalidSlug
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		logger.Error("failed to hash password", logger.Err(err))

		return nil, err
	}

	user := model.NewUser(req.Email, req.Username, req.Name, req.Phone, hashedPassword)
	tenant := model.NewTenant(req.TenantName, req.TenantSlug)
	membership := model.NewTenantMembership(tenant.ID, user.ID, model.TenantRoleOwner)

	if err := s.tenantRepo.CreateWithOwner(ctx, user, tenant, membership); err != nil {
		if serrors.Is(err, tenantrepo.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}

		if serrors.Is(err, tenantrepo.ErrTenantAlreadyExists) {
			return nil, ErrTenantAlreadyExists
		}

		logger.Error("failed to create user and tenant", logger.Err(err))

		return nil, err
	}

	return &dto.RegisterResponse{
		User: *s.userToResponse(user),
		Tenant: dto.TenantSummary{
			ID:   tenant.ID,
			Name: tenant.Name,
			Slug: tenant.Slug,
			Role: membership.Role,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest, deviceInfo dto.DeviceInfo) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		logger.Error("failed to get user by email", logger.Err(err))

		return nil, err
	}

	if !ComparePassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.issueSession(ctx, user, deviceInfo)
}

func (s *authService) RefreshToken(ctx context.Context, refreshTokenString string) (*dto.TokenResponse, error) {
	claims, err := s.tokenManager.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != jwt.RefreshToken {
		return nil, ErrInvalidToken
	}

	if claims.ID == "" {
		return nil, ErrInvalidToken
	}

	session, err := s.sessionRepo.GetSessionByRID(ctx, claims.ID)
	if err != nil {
		if serrors.Is(err, sessionrepo.ErrSessionNotFound) {
			return nil, ErrSessionNotFound
		}

		logger.Error("failed to get session by RID", logger.Err(err))

		return nil, err
	}

	if session.Revoked {
		return nil, ErrSessionRevoked
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	if session.UserID != claims.UserID {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		logger.Error("failed to get user by ID", logger.Err(err))

		return nil, err
	}

	newRID := uuid.NewString()
	if err := s.sessionRepo.UpdateSessionRID(ctx, session.ID, newRID); err != nil {
		logger.Error("failed to update session RID", logger.Err(err))

		return nil, err
	}

	return s.issueTokens(user, session.ID, newRID)
}

func (s *authService) ValidateSession(ctx context.Context, sessionID uuid.UUID) (*model.Session, error) {
	session, err := s.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		if serrors.Is(err, sessionrepo.ErrSessionNotFound) {
			return nil, ErrSessionNotFound
		}

		return nil, err
	}

	if session.Revoked {
		return nil, ErrSessionRevoked
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return session, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *authService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.sessionRepo.RevokeSession(ctx, sessionID); err != nil {
		logger.Error("failed to revoke session", logger.Err(err))

		return err
	}

	return nil
}

func (s *authService) CurrentSession(ctx context.Context, userID uuid.UUID) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		logger.Error("failed to get user by ID", logger.Err(err))

		return nil, err
	}

	return &dto.AuthResponse{
		User: *s.userToResponse(user),
	}, nil
}
