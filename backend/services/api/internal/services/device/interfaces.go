package device

import (
	"context"

	"ava/api/internal/dto"
	devicerepo "ava/api/internal/repository/device"
	eventsvc "ava/api/internal/services/event"
	"ava/pkg/wire"

	"github.com/google/uuid"
)

type Service interface {
	SyncFromHub(ctx context.Context, tenantID, hubID uuid.UUID, req *dto.SyncDevicesRequest) ([]*dto.DeviceResponse, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.DeviceResponse, error)
	ListByHub(ctx context.Context, tenantID, hubID uuid.UUID) ([]*dto.DeviceResponse, error)
	Update(ctx context.Context, tenantID, deviceID uuid.UUID, req *dto.UpdateDeviceRequest) (*dto.DeviceResponse, error)
	SendCommand(ctx context.Context, tenantID, deviceID uuid.UUID, req *dto.SendCommandRequest) (*dto.CommandAcceptedResponse, error)
	Apply(ctx context.Context, tenantID uuid.UUID, req *dto.ApplyRequest) (*dto.ApplyResponse, error)
	ApplyReportedState(ctx context.Context, hubID uuid.UUID, externalID string, state wire.State) error
	MarkHubOffline(ctx context.Context, tenantID, hubID uuid.UUID) error
}

type Commander interface {
	Publish(ctx context.Context, topic string, payload []byte, retained bool) error
}

type deviceService struct {
	deviceRepo devicerepo.Repository
	commander  Commander
	events     eventsvc.Service
}

func NewService(deviceRepo devicerepo.Repository, commander Commander, events eventsvc.Service) Service {
	return &deviceService{
		deviceRepo: deviceRepo,
		commander:  commander,
		events:     events,
	}
}
