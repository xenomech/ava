package tenant

import (
	"context"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, tenant *model.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*model.Tenant, error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Tenant, error)
	Update(ctx context.Context, tenant *model.Tenant) error
	CreateWithOwner(ctx context.Context, user *model.User, tenant *model.Tenant, membership *model.TenantMembership) error
	CreateWithMembership(ctx context.Context, tenant *model.Tenant, membership *model.TenantMembership) error
}

type tenantRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &tenantRepository{db: db}
}
