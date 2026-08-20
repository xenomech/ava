package hub

import (
	"context"
	"errors"
	"time"

	"ava/api/internal/model"
	"ava/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *hubRepository) CreateAuthorization(ctx context.Context, auth *model.HubAuthorization) error {
	return r.db.WithContext(ctx).Create(auth).Error
}

func (r *hubRepository) GetAuthorizationByDeviceCode(ctx context.Context, deviceCode string) (*model.HubAuthorization, error) {
	var auth model.HubAuthorization

	err := r.db.WithContext(ctx).Where("device_code = ?", deviceCode).First(&auth).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAuthorizationNotFound
	}

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *hubRepository) GetAuthorizationByUserCode(ctx context.Context, userCode string) (*model.HubAuthorization, error) {
	var auth model.HubAuthorization

	err := r.db.WithContext(ctx).Where("user_code = ?", userCode).First(&auth).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAuthorizationNotFound
	}

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *hubRepository) UpdateAuthorization(ctx context.Context, auth *model.HubAuthorization) error {
	result := r.db.WithContext(ctx).Model(&model.HubAuthorization{}).Where("id = ?", auth.ID).Updates(map[string]any{
		"status":      auth.Status,
		"tenant_id":   auth.TenantID,
		"user_id":     auth.UserID,
		"hub_id":      auth.HubID,
		"approved_at": auth.ApprovedAt,
		"updated_at":  time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrAuthorizationNotFound
	}

	return nil
}

func (r *hubRepository) TouchAuthorizationPoll(ctx context.Context, id uuid.UUID, at time.Time) error {
	return r.db.WithContext(ctx).Model(&model.HubAuthorization{}).Where("id = ?", id).
		Update("last_polled_at", at).Error
}

func (r *hubRepository) ApproveWithDevice(ctx context.Context, auth *model.HubAuthorization, hub *model.Hub) error {
	return r.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Create(hub).Error; err != nil {
			return err
		}

		auth.HubID = &hub.ID

		result := dbTx.Model(&model.HubAuthorization{}).Where("id = ? AND status = ?", auth.ID, model.HubAuthStatusPending).
			Updates(map[string]any{
				"status":      auth.Status,
				"tenant_id":   auth.TenantID,
				"user_id":     auth.UserID,
				"hub_id":      auth.HubID,
				"approved_at": auth.ApprovedAt,
				"updated_at":  time.Now(),
			})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrAuthorizationNotFound
		}

		logger.Info("hubRepository.ApproveWithDevice",
			logger.Any("hub.ID", hub.ID),
			logger.Any("tenant.ID", hub.TenantID),
		)

		return nil
	})
}

func (r *hubRepository) Create(ctx context.Context, hub *model.Hub) error {
	return r.db.WithContext(ctx).Create(hub).Error
}

func (r *hubRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*model.Hub, error) {
	var hub model.Hub

	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&hub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrHubNotFound
	}

	if err != nil {
		return nil, err
	}

	return &hub, nil
}

func (r *hubRepository) GetByIDUnscoped(ctx context.Context, id uuid.UUID) (*model.Hub, error) {
	var hub model.Hub

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&hub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrHubNotFound
	}

	if err != nil {
		return nil, err
	}

	return &hub, nil
}

func (r *hubRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Hub, error) {
	var hub model.Hub

	err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&hub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrHubNotFound
	}

	if err != nil {
		return nil, err
	}

	return &hub, nil
}

func (r *hubRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*model.Hub, error) {
	var hubs []*model.Hub

	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at ASC").Find(&hubs).Error
	if err != nil {
		return nil, err
	}

	return hubs, nil
}

func (r *hubRepository) UpdateRefreshToken(ctx context.Context, hubID uuid.UUID, refreshToken string) error {
	result := r.db.WithContext(ctx).Model(&model.Hub{}).Where("id = ?", hubID).Updates(map[string]any{
		"refresh_token": refreshToken,
		"updated_at":    time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrHubNotFound
	}

	return nil
}

func (r *hubRepository) TouchLastSeen(ctx context.Context, hubID uuid.UUID, at time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Hub{}).Where("id = ?", hubID).
		Update("last_seen_at", at).Error
}

func (r *hubRepository) Revoke(ctx context.Context, tenantID, hubID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&model.Hub{}).Where("tenant_id = ? AND id = ?", tenantID, hubID).
		Updates(map[string]any{
			"status":     model.HubStatusRevoked,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrHubNotFound
	}

	return nil
}
