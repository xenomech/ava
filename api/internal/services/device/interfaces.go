package device

import (
	"context"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	devicerepo "ava/api/internal/repository/device"
	tenantrepo "ava/api/internal/repository/tenant"
	"ava/api/internal/services/auth/jwt"

	"github.com/google/uuid"
)

type Service interface {
	RequestCode(ctx context.Context, req *dto.DeviceCodeRequest) (*dto.DeviceCodeResponse, error)
	Poll(ctx context.Context, deviceCode string) (*dto.DeviceTokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.DeviceTokenResponse, error)
	Activate(ctx context.Context, tenantID, userID uuid.UUID, userCode string) (*dto.DeviceResponse, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.DeviceResponse, error)
	Revoke(ctx context.Context, tenantID, deviceID uuid.UUID) error
	Heartbeat(ctx context.Context, deviceID uuid.UUID) error
	ValidateDevice(ctx context.Context, deviceID uuid.UUID) (*model.Device, error)
}

type deviceService struct {
	deviceRepo   devicerepo.Repository
	tenantRepo   tenantrepo.Repository
	tokenManager jwt.TokenManager
}

func NewService(deviceRepo devicerepo.Repository, tenantRepo tenantrepo.Repository) Service {
	return &deviceService{
		deviceRepo:   deviceRepo,
		tenantRepo:   tenantRepo,
		tokenManager: jwt.NewTokenManager(),
	}
}
