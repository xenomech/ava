package device

import (
	"context"
	"errors"
	"time"

	"ava/api/internal/model"
	"ava/api/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *deviceRepository) CreateAuthorization(ctx context.Context, auth *model.DeviceAuthorization) error {
	return r.db.WithContext(ctx).Create(auth).Error
}

func (r *deviceRepository) GetAuthorizationByDeviceCode(ctx context.Context, deviceCode string) (*model.DeviceAuthorization, error) {
	var auth model.DeviceAuthorization

	err := r.db.WithContext(ctx).Where("device_code = ?", deviceCode).First(&auth).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAuthorizationNotFound
	}

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *deviceRepository) GetAuthorizationByUserCode(ctx context.Context, userCode string) (*model.DeviceAuthorization, error) {
	var auth model.DeviceAuthorization

	err := r.db.WithContext(ctx).Where("user_code = ?", userCode).First(&auth).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAuthorizationNotFound
	}

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *deviceRepository) UpdateAuthorization(ctx context.Context, auth *model.DeviceAuthorization) error {
	result := r.db.WithContext(ctx).Model(&model.DeviceAuthorization{}).Where("id = ?", auth.ID).Updates(map[string]any{
		"status":      auth.Status,
		"tenant_id":   auth.TenantID,
		"user_id":     auth.UserID,
		"device_id":   auth.DeviceID,
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

func (r *deviceRepository) TouchAuthorizationPoll(ctx context.Context, id uuid.UUID, at time.Time) error {
	return r.db.WithContext(ctx).Model(&model.DeviceAuthorization{}).Where("id = ?", id).
		Update("last_polled_at", at).Error
}

func (r *deviceRepository) ApproveWithDevice(ctx context.Context, auth *model.DeviceAuthorization, device *model.Device) error {
	return r.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Create(device).Error; err != nil {
			return err
		}

		auth.DeviceID = &device.ID

		result := dbTx.Model(&model.DeviceAuthorization{}).Where("id = ? AND status = ?", auth.ID, model.DeviceAuthStatusPending).
			Updates(map[string]any{
				"status":      auth.Status,
				"tenant_id":   auth.TenantID,
				"user_id":     auth.UserID,
				"device_id":   auth.DeviceID,
				"approved_at": auth.ApprovedAt,
				"updated_at":  time.Now(),
			})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrAuthorizationNotFound
		}

		logger.Info("deviceRepository.ApproveWithDevice",
			logger.Any("device.ID", device.ID),
			logger.Any("tenant.ID", device.TenantID),
		)

		return nil
	})
}

func (r *deviceRepository) Create(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

func (r *deviceRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*model.Device, error) {
	var device model.Device

	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDeviceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *deviceRepository) GetByIDUnscoped(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	var device model.Device

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDeviceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *deviceRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Device, error) {
	var device model.Device

	err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDeviceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *deviceRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*model.Device, error) {
	var devices []*model.Device

	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("created_at ASC").Find(&devices).Error
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *deviceRepository) UpdateRefreshToken(ctx context.Context, deviceID uuid.UUID, refreshToken string) error {
	result := r.db.WithContext(ctx).Model(&model.Device{}).Where("id = ?", deviceID).Updates(map[string]any{
		"refresh_token": refreshToken,
		"updated_at":    time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}

	return nil
}

func (r *deviceRepository) TouchLastSeen(ctx context.Context, deviceID uuid.UUID, at time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Device{}).Where("id = ?", deviceID).
		Update("last_seen_at", at).Error
}

func (r *deviceRepository) Revoke(ctx context.Context, tenantID, deviceID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&model.Device{}).Where("tenant_id = ? AND id = ?", tenantID, deviceID).
		Updates(map[string]any{
			"status":     model.DeviceStatusRevoked,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}

	return nil
}
