package membership

import (
	"context"

	"ava/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, membership *model.TenantMembership) error
	GetByTenantAndUser(ctx context.Context, tenantID, userID uuid.UUID) (*model.TenantMembership, error)
	GetByInviteToken(ctx context.Context, inviteToken string) (*model.TenantMembership, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*model.TenantMembership, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.TenantMembership, error)
	CountByRole(ctx context.Context, tenantID uuid.UUID, role model.TenantRole) (int64, error)
	UpdateRole(ctx context.Context, tenantID, userID uuid.UUID, role model.TenantRole) error
	Activate(ctx context.Context, tenantID, userID uuid.UUID) error
	Delete(ctx context.Context, tenantID, userID uuid.UUID) error
}

type membershipRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &membershipRepository{db: db}
}
