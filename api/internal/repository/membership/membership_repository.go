package membership

import (
	"context"
	"errors"
	"time"

	"ava/api/internal/model"
	"ava/api/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *membershipRepository) Create(ctx context.Context, membership *model.TenantMembership) error {
	err := r.db.WithContext(ctx).Create(membership).Error

	logger.Info("membershipRepository.Create",
		logger.Any("membership.ID", membership.ID),
		logger.Any("error", err),
	)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrMembershipAlreadyExists
	}

	return err
}

func (r *membershipRepository) GetByTenantAndUser(ctx context.Context, tenantID, userID uuid.UUID) (*model.TenantMembership, error) {
	var membership model.TenantMembership

	err := r.db.WithContext(ctx).Where("tenant_id = ? AND user_id = ?", tenantID, userID).First(&membership).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrMembershipNotFound
	}

	if err != nil {
		return nil, err
	}

	return &membership, nil
}

func (r *membershipRepository) GetByInviteToken(ctx context.Context, inviteToken string) (*model.TenantMembership, error) {
	var membership model.TenantMembership

	err := r.db.WithContext(ctx).
		Where("invite_token = ? AND status = ? AND invite_expires_at > ?", inviteToken, model.MembershipStatusInvited, time.Now()).
		First(&membership).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrMembershipNotFound
	}

	if err != nil {
		return nil, err
	}

	return &membership, nil
}

func (r *membershipRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*model.TenantMembership, error) {
	var memberships []*model.TenantMembership

	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at ASC").Find(&memberships).Error
	if err != nil {
		return nil, err
	}

	return memberships, nil
}

func (r *membershipRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.TenantMembership, error) {
	var memberships []*model.TenantMembership

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at ASC").Find(&memberships).Error
	if err != nil {
		return nil, err
	}

	return memberships, nil
}

func (r *membershipRepository) CountByRole(ctx context.Context, tenantID uuid.UUID, role model.TenantRole) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&model.TenantMembership{}).
		Where("tenant_id = ? AND role = ? AND status = ?", tenantID, role, model.MembershipStatusActive).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *membershipRepository) UpdateRole(ctx context.Context, tenantID, userID uuid.UUID, role model.TenantRole) error {
	result := r.db.WithContext(ctx).Model(&model.TenantMembership{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Updates(map[string]any{
			"role":       role,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrMembershipNotFound
	}

	return nil
}

func (r *membershipRepository) Activate(ctx context.Context, tenantID, userID uuid.UUID) error {
	now := time.Now()

	result := r.db.WithContext(ctx).Model(&model.TenantMembership{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Updates(map[string]any{
			"status":            model.MembershipStatusActive,
			"invite_token":      "",
			"invite_expires_at": nil,
			"joined_at":         &now,
			"updated_at":        now,
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrMembershipNotFound
	}

	return nil
}

func (r *membershipRepository) Delete(ctx context.Context, tenantID, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Unscoped().
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Delete(&model.TenantMembership{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrMembershipNotFound
	}

	return nil
}
