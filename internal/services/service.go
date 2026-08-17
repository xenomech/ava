package services

import (
	"ava/internal/repository"
	authsvc "ava/internal/services/auth"
	healthsvc "ava/internal/services/health"
)

type Service struct {
	Auth   authsvc.Service
	Health healthsvc.Service
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth:   authsvc.NewService(repo.User, repo.Tenant, repo.Membership, repo.Session),
		Health: healthsvc.NewService(),
	}
}
