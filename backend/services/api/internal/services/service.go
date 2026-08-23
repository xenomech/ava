package services

import (
	"ava/api/internal/repository"
	authsvc "ava/api/internal/services/auth"
	devicesvc "ava/api/internal/services/device"
	eventsvc "ava/api/internal/services/event"
	flowsvc "ava/api/internal/services/flow"
	healthsvc "ava/api/internal/services/health"
	hubsvc "ava/api/internal/services/hub"
	roomsvc "ava/api/internal/services/room"
	tenantsvc "ava/api/internal/services/tenant"
)

type Commander = devicesvc.Commander

type Provisioner = hubsvc.BrokerProvisioner

type Service struct {
	Auth   authsvc.Service
	Tenant tenantsvc.Service
	Flow   flowsvc.Service
	Hub    hubsvc.Service
	Room   roomsvc.Service
	Device devicesvc.Service
	Event  eventsvc.Service
	Health healthsvc.Service
}

func NewService(repo *repository.Repository, commander devicesvc.Commander, provisioner Provisioner) *Service {
	tenantService := tenantsvc.NewService(repo.Tenant, repo.Membership, repo.User, repo.Session)
	eventService := eventsvc.NewService()
	deviceService := devicesvc.NewService(repo.Device, commander, eventService)

	return &Service{
		Auth:   authsvc.NewService(repo.User, repo.Tenant, repo.Membership, repo.Session, repo.Token),
		Tenant: tenantService,
		Flow:   flowsvc.NewService(repo.Flow, tenantService, repo.Membership, repo.Hub),
		Room:   roomsvc.NewService(repo.Room),
		Hub:    hubsvc.NewService(repo.Hub, repo.Tenant, eventService, deviceService, provisioner),
		Device: deviceService,
		Event:  eventService,
		Health: healthsvc.NewService(),
	}
}
