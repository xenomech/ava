package services

import (
	"ava/internal/repository"
	authsvc "ava/internal/services/auth"
	flowsvc "ava/internal/services/flow"
	healthsvc "ava/internal/services/health"
	tenantsvc "ava/internal/services/tenant"
)

type Service struct {
	Auth   authsvc.Service
	Tenant tenantsvc.Service
	Flow   flowsvc.Service
	Health healthsvc.Service
}

func NewService(repo *repository.Repository) *Service {
	tenantService := tenantsvc.NewService(repo.Tenant, repo.Membership, repo.User, repo.Session)

	return &Service{
		Auth:   authsvc.NewService(repo.User, repo.Tenant, repo.Membership, repo.Session, repo.Token),
		Tenant: tenantService,
		Flow:   flowsvc.NewService(repo.Flow, tenantService, repo.User, repo.Membership),
		Health: healthsvc.NewService(),
	}
}
