package hub

import (
	"context"
	"time"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateAuthorization(ctx context.Context, auth *model.HubAuthorization) error
	GetAuthorizationByDeviceCode(ctx context.Context, deviceCode string) (*model.HubAuthorization, error)
	GetAuthorizationByUserCode(ctx context.Context, userCode string) (*model.HubAuthorization, error)
	UpdateAuthorization(ctx context.Context, auth *model.HubAuthorization) error
	TouchAuthorizationPoll(ctx context.Context, id uuid.UUID, at time.Time) error
	ApproveWithDevice(ctx context.Context, auth *model.HubAuthorization, hub *model.Hub) error

	Create(ctx context.Context, hub *model.Hub) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*model.Hub, error)
	GetByIDUnscoped(ctx context.Context, id uuid.UUID) (*model.Hub, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Hub, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*model.Hub, error)
	UpdateRefreshToken(ctx context.Context, hubID uuid.UUID, refreshToken string) error
	TouchLastSeen(ctx context.Context, hubID uuid.UUID, at time.Time) error
	Revoke(ctx context.Context, tenantID, hubID uuid.UUID) error
}

type hubRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &hubRepository{db: db}
}
