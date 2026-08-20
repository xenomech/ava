package device

import (
	"context"
	"encoding/json"

	"ava/api/internal/dto"
	devicerepo "ava/api/internal/repository/device"
	hubrepo "ava/api/internal/repository/hub"
	tenantrepo "ava/api/internal/repository/tenant"
	eventsvc "ava/api/internal/services/event"

	"github.com/google/uuid"
)

type Service interface {
	SyncFromHub(ctx context.Context, tenantID, hubID uuid.UUID, req *dto.SyncDevicesRequest) ([]*dto.DeviceResponse, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.DeviceResponse, error)
	ListByHub(ctx context.Context, tenantID, hubID uuid.UUID) ([]*dto.DeviceResponse, error)
	Update(ctx context.Context, tenantID, deviceID uuid.UUID, req *dto.UpdateDeviceRequest) (*dto.DeviceResponse, error)
	SendCommand(ctx context.Context, tenantID, deviceID uuid.UUID, req *dto.SendCommandRequest) (*dto.CommandAcceptedResponse, error)
	ApplyReportedState(ctx context.Context, hubID uuid.UUID, externalID string, state json.RawMessage) error
	MarkHubOffline(ctx context.Context, tenantID, hubID uuid.UUID) error
}

type Commander interface {
	Publish(ctx context.Context, topic string, payload []byte, retained bool) error
}

type deviceService struct {
	deviceRepo devicerepo.Repository
	hubRepo    hubrepo.Repository
	tenantRepo tenantrepo.Repository
	commander  Commander
	events     eventsvc.Service
}

func NewService(
	deviceRepo devicerepo.Repository,
	hubRepo hubrepo.Repository,
	tenantRepo tenantrepo.Repository,
	commander Commander,
	events eventsvc.Service,
) Service {
	return &deviceService{
		deviceRepo: deviceRepo,
		hubRepo:    hubRepo,
		tenantRepo: tenantRepo,
		commander:  commander,
		events:     events,
	}
}
