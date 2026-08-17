package services

import (
	"ava/internal/repository"
	authsvc "ava/internal/services/auth"
	healthsvc "ava/internal/services/health"
	tenantsvc "ava/internal/services/tenant"
)

type Service struct {
	Auth   authsvc.Service
	Tenant tenantsvc.Service
	Health healthsvc.Service
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth:   authsvc.NewService(repo.User, repo.Tenant, repo.Membership, repo.Session),
		Tenant: tenantsvc.NewService(repo.Tenant, repo.Membership, repo.User, repo.Session),
		Health: healthsvc.NewService(),
	}
}
