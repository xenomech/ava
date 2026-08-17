package auth

import (
	"context"

	"ava/internal/dto"
	sessionrepo "ava/internal/repository/session"
	userrepo "ava/internal/repository/user"
	"ava/internal/services/auth/jwt"
)

type Service interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
}

type authService struct {
	userRepo     userrepo.Repository
	sessionRepo  sessionrepo.Repository
	tokenManager jwt.TokenManager
}

func NewService(
	userRepo userrepo.Repository,
	sessionRepo sessionrepo.Repository,
) Service {
	return &authService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenManager: jwt.NewTokenManager(),
	}
}
