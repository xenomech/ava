package device

import (
	"context"

	"ava/api/internal/dto"
	devicerepo "ava/api/internal/repository/device"
	tenantrepo "ava/api/internal/repository/tenant"

	"github.com/google/uuid"
)

type Service interface {
	SyncFromHub(ctx context.Context, tenantID, hubID uuid.UUID, req *dto.SyncDevicesRequest) ([]*dto.DeviceResponse, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.DeviceResponse, error)
	ListByHub(ctx context.Context, tenantID, hubID uuid.UUID) ([]*dto.DeviceResponse, error)
	Rename(ctx context.Context, tenantID, deviceID uuid.UUID, name string) (*dto.DeviceResponse, error)
	SendCommand(ctx context.Context, tenantID, deviceID uuid.UUID, req *dto.SendCommandRequest) (*dto.CommandAcceptedResponse, error)
}

type Commander interface {
	Publish(ctx context.Context, topic string, payload []byte, retained bool) error
}

type deviceService struct {
	deviceRepo devicerepo.Repository
	tenantRepo tenantrepo.Repository
	commander  Commander
}

func NewService(deviceRepo devicerepo.Repository, tenantRepo tenantrepo.Repository, commander Commander) Service {
	return &deviceService{
		deviceRepo: deviceRepo,
		tenantRepo: tenantRepo,
		commander:  commander,
	}
}
