package device

import (
	"context"

	"ava/api/internal/dto"
	devicerepo "ava/api/internal/repository/device"

	"github.com/google/uuid"
)

type Service interface {
	SyncFromHub(ctx context.Context, tenantID, hubID uuid.UUID, req *dto.SyncDevicesRequest) ([]*dto.DeviceResponse, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.DeviceResponse, error)
	ListByHub(ctx context.Context, tenantID, hubID uuid.UUID) ([]*dto.DeviceResponse, error)
	Rename(ctx context.Context, tenantID, deviceID uuid.UUID, name string) (*dto.DeviceResponse, error)
}

type deviceService struct {
	deviceRepo devicerepo.Repository
}

func NewService(deviceRepo devicerepo.Repository) Service {
	return &deviceService{deviceRepo: deviceRepo}
}
