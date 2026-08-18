package services

import (
	"ava/api/internal/repository"
	authsvc "ava/api/internal/services/auth"
	devicesvc "ava/api/internal/services/device"
	flowsvc "ava/api/internal/services/flow"
	healthsvc "ava/api/internal/services/health"
	tenantsvc "ava/api/internal/services/tenant"
)

type Service struct {
	Auth   authsvc.Service
	Tenant tenantsvc.Service
	Flow   flowsvc.Service
	Device devicesvc.Service
	Health healthsvc.Service
}

func NewService(repo *repository.Repository) *Service {
	tenantService := tenantsvc.NewService(repo.Tenant, repo.Membership, repo.User, repo.Session)

	return &Service{
		Auth:   authsvc.NewService(repo.User, repo.Tenant, repo.Membership, repo.Session, repo.Token),
		Tenant: tenantService,
		Flow:   flowsvc.NewService(repo.Flow, tenantService, repo.User, repo.Membership),
		Device: devicesvc.NewService(repo.Device, repo.Tenant),
		Health: healthsvc.NewService(),
	}
}
