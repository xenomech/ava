package auth

import (
	"context"
	"time"

	"ava/internal/dto"
	"ava/internal/model"
	sessionrepo "ava/internal/repository/session"
	userrepo "ava/internal/repository/user"
	"ava/internal/services/auth/jwt"
	"ava/pkg/logger"
	"ava/pkg/serrors"

	"github.com/google/uuid"
)

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		logger.Error("failed to hash password", logger.Err(err))

		return nil, err
	}

	user := model.NewUser(req.Email, req.Username, req.Name, req.Phone, hashedPassword)

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		if serrors.Is(err, userrepo.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}

		logger.Error("failed to create user", logger.Err(err))

		return nil, err
	}

	return &dto.RegisterResponse{
		User: *s.userToResponse(user),
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
