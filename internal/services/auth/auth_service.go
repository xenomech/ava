package auth

import (
	"context"

	"ava/internal/dto"
	"ava/internal/model"
	userrepo "ava/internal/repository/user"
	"ava/pkg/logger"
	"ava/pkg/serrors"
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
