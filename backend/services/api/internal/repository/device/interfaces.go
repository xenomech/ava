package device

import (
	"context"
	"encoding/json"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	SyncHubDevices(ctx context.Context, tenantID, hubID uuid.UUID, devices []*model.Device) error
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*model.Device, error)
	ListByHub(ctx context.Context, tenantID, hubID uuid.UUID) ([]*model.Device, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*model.Device, error)
	Update(ctx context.Context, tenantID, id uuid.UUID, fields map[string]any) error
	ApplyState(ctx context.Context, hubID uuid.UUID, externalID string, state json.RawMessage) (*model.Device, error)
}

type deviceRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &deviceRepository{db: db}
}
