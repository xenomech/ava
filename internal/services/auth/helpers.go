package auth

import (
	"context"
	"time"

	"ava/internal/dto"
	"ava/internal/model"
	"ava/internal/services/auth/jwt"
	"ava/pkg/logger"

	"github.com/google/uuid"
)

func (s *authService) issueSession(ctx context.Context, user *model.User, deviceInfo dto.DeviceInfo) (*dto.AuthResponse, error) {
	rid := uuid.NewString()

	session := model.NewSession(
		user.ID,
		deviceInfo.DeviceName,
		deviceInfo.IPAddress,
		deviceInfo.UserAgent,
		rid,
		time.Now().Add(s.tokenManager.GetRefreshExpiry()),
	)

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		logger.Error("failed to create session", logger.Err(err))

		return nil, err
	}

	tokens, err := s.issueTokens(user, session.ID, rid)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User:   *s.userToResponse(user),
		Tokens: tokens,
	}, nil
}

func (s *authService) issueTokens(user *model.User, sessionID uuid.UUID, rid string) (*dto.TokenResponse, error) {
	accessToken, err := s.tokenManager.GenerateToken(user, sessionID, jwt.AccessToken, "")
	if err != nil {
		logger.Error("failed to generate access token", logger.Err(err))

		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateToken(user, sessionID, jwt.RefreshToken, rid)
	if err != nil {
		logger.Error("failed to generate refresh token", logger.Err(err))

		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenManager.GetAccessExpiry().Seconds()),
	}, nil
}

func (s *authService) userToResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		Name:            user.Name,
		Phone:           user.Phone,
		EmailVerified:   user.EmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:       user.CreatedAt,
	}
}
