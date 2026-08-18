package device

import (
	"context"
	"time"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateAuthorization(ctx context.Context, auth *model.DeviceAuthorization) error
	GetAuthorizationByDeviceCode(ctx context.Context, deviceCode string) (*model.DeviceAuthorization, error)
	GetAuthorizationByUserCode(ctx context.Context, userCode string) (*model.DeviceAuthorization, error)
	UpdateAuthorization(ctx context.Context, auth *model.DeviceAuthorization) error
	TouchAuthorizationPoll(ctx context.Context, id uuid.UUID, at time.Time) error
	ApproveWithDevice(ctx context.Context, auth *model.DeviceAuthorization, device *model.Device) error

	Create(ctx context.Context, device *model.Device) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*model.Device, error)
	GetByIDUnscoped(ctx context.Context, id uuid.UUID) (*model.Device, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Device, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*model.Device, error)
	UpdateRefreshToken(ctx context.Context, deviceID uuid.UUID, refreshToken string) error
	TouchLastSeen(ctx context.Context, deviceID uuid.UUID, at time.Time) error
	Revoke(ctx context.Context, tenantID, deviceID uuid.UUID) error
}

type deviceRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &deviceRepository{db: db}
}
