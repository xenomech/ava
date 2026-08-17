package tenant

import (
	"context"
	"errors"

	"ava/internal/model"
	"ava/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *tenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	err := r.db.WithContext(ctx).Create(tenant).Error

	logger.Info("tenantRepository.Create", logger.Any("tenant.ID", tenant.ID), logger.Any("error", err))

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrTenantAlreadyExists
	}

	return err
}

func (r *tenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Tenant, error) {
	var tenant model.Tenant

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tenant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTenantNotFound
	}

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (r *tenantRepository) GetBySlug(ctx context.Context, slug string) (*model.Tenant, error) {
	var tenant model.Tenant

	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tenant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTenantNotFound
	}

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (r *tenantRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Tenant, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var tenants []*model.Tenant

	err := r.db.WithContext(ctx).Where("id IN ?", ids).Order("created_at ASC").Find(&tenants).Error
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

func (r *tenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	result := r.db.WithContext(ctx).Model(tenant).Where("id = ?", tenant.ID).Updates(tenant)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTenantNotFound
	}

	return nil
}

func (r *tenantRepository) CreateWithMembership(ctx context.Context, tenant *model.Tenant, membership *model.TenantMembership) error {
	return r.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Create(tenant).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrTenantAlreadyExists
			}

			return err
		}

		membership.TenantID = tenant.ID

		if err := dbTx.Create(membership).Error; err != nil {
			return err
		}

		logger.Info("tenantRepository.CreateWithMembership",
			logger.Any("tenant.ID", tenant.ID),
			logger.Any("user.ID", membership.UserID),
		)

		return nil
	})
}

func (r *tenantRepository) CreateWithOwner(ctx context.Context, user *model.User, tenant *model.Tenant, membership *model.TenantMembership) error {
	return r.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Create(user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrUserAlreadyExists
			}

			return err
		}

		if err := dbTx.Create(tenant).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrTenantAlreadyExists
			}

			return err
		}

		membership.UserID = user.ID
		membership.TenantID = tenant.ID

		if err := dbTx.Create(membership).Error; err != nil {
			return err
		}

		logger.Info("tenantRepository.CreateWithOwner",
			logger.Any("user.ID", user.ID),
			logger.Any("tenant.ID", tenant.ID),
		)

		return nil
	})
}
