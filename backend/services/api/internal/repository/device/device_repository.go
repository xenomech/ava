package device

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"ava/api/internal/model"
	"ava/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// What a sweep is allowed to overwrite on a device it has seen before.
//
// Not "state". A sweep describes a device as it was when discovery read it, and
// discovery takes several seconds — so writing that reading over the stored one
// let a six-second-old snapshot undo a change made a second ago. The bulb had
// the new colour and the app showed the old white, which is the worst way round
// for the two to disagree.
//
// State arrives over the broker instead, where it is merged rather than
// replaced: a reading can correct a trait or retire it, but it cannot revert
// something newer than itself. A device seen for the first time still gets its
// state at insert, because there is nothing there to protect.
var hubOwnedColumns = []string{
	"kind",
	"vendor",
	"model",
	"ip",
	"parent",
	"capabilities",
	"status",
	"last_seen_at",
	"updated_at",
}

func (r *deviceRepository) SyncHubDevices(ctx context.Context, tenantID, hubID uuid.UUID, devices []*model.Device) error {
	return r.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		seen := make([]string, 0, len(devices))

		for _, device := range devices {
			device.TenantID = tenantID
			device.HubID = hubID
			seen = append(seen, device.ExternalID)
		}

		if len(devices) > 0 {
			err := dbTx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "hub_id"}, {Name: "external_id"}},
				DoUpdates: clause.AssignmentColumns(hubOwnedColumns),
			}).Create(&devices).Error
			if err != nil {
				return err
			}
		}

		query := dbTx.Model(&model.Device{}).Where("tenant_id = ? AND hub_id = ?", tenantID, hubID)
		if len(seen) > 0 {
			query = query.Where("external_id NOT IN ?", seen)
		}

		result := query.Updates(map[string]any{
			"status":     model.DeviceStatusOffline,
			"updated_at": time.Now(),
		})
		if result.Error != nil {
			return result.Error
		}

		logger.Info("deviceRepository.SyncHubDevices",
			logger.Any("hub.ID", hubID),
			logger.Int("reported", len(devices)),
			logger.Int("marked_offline", int(result.RowsAffected)),
		)

		return nil
	})
}

func (r *deviceRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*model.Device, error) {
	var devices []*model.Device

	err := r.db.WithContext(ctx).Preload("Room").
		Where("tenant_id = ?", tenantID).Order("created_at ASC").Find(&devices).Error
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *deviceRepository) ListByHub(ctx context.Context, tenantID, hubID uuid.UUID) ([]*model.Device, error) {
	var devices []*model.Device

	err := r.db.WithContext(ctx).Preload("Room").
		Where("tenant_id = ? AND hub_id = ?", tenantID, hubID).
		Order("created_at ASC").Find(&devices).Error
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *deviceRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*model.Device, error) {
	var device model.Device

	err := r.db.WithContext(ctx).Preload("Room").
		Where("tenant_id = ? AND id = ?", tenantID, id).First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDeviceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *deviceRepository) ListWithRelations(
	ctx context.Context,
	tenantID uuid.UUID,
	ids []uuid.UUID,
) ([]*model.Device, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var devices []*model.Device

	err := r.db.WithContext(ctx).
		Preload("Hub").
		Preload("Tenant").
		Preload("Room").
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Find(&devices).Error
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *deviceRepository) GetWithRelations(ctx context.Context, tenantID, id uuid.UUID) (*model.Device, error) {
	var device model.Device

	err := r.db.WithContext(ctx).
		Joins("Hub").
		Joins("Tenant").
		Preload("Room").
		Where("devices.tenant_id = ? AND devices.id = ?", tenantID, id).
		First(&device).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDeviceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *deviceRepository) Update(ctx context.Context, tenantID, id uuid.UUID, fields map[string]any) error {
	fields["updated_at"] = time.Now()

	result := r.db.WithContext(ctx).Model(&model.Device{}).Where("tenant_id = ? AND id = ?", tenantID, id).
		Updates(fields)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}

	return nil
}

func (r *deviceRepository) ApplyState(
	ctx context.Context,
	hubID uuid.UUID,
	externalID string,
	state json.RawMessage,
	cleared []string,
) (*model.Device, error) {
	var updated model.Device

	/* Merge what was reported, then remove the traits reported as null. Built
	   one `- ?` at a time rather than as an array parameter so it does not
	   depend on the driver encoding a text[]; a device reports a handful of
	   traits, never enough for the difference to matter. */
	/* Parenthesised, because `-` binds tighter than `||`: written as
	   `state || patch - 'color'` Postgres removes the key from the patch and
	   then merges, which does nothing at all. The retirement read correctly and
	   had never once run, so a bulb moved from colour to white kept its old
	   colour key alongside the new temperature — and the app, seeing a
	   temperature, drew white over the colour it had just been asked for. */
	merge := "(coalesce(state, '{}'::jsonb) || ?::jsonb)"
	args := []any{string(state)}

	for _, trait := range cleared {
		merge += " - ?::text"

		args = append(args, trait)
	}

	result := r.db.WithContext(ctx).Model(&updated).
		Clauses(clause.Returning{}).
		Where("hub_id = ? AND external_id = ?", hubID, externalID).
		Updates(map[string]any{
			"state":        gorm.Expr(merge, args...),
			"status":       model.DeviceStatusOnline,
			"last_seen_at": time.Now(),
			"updated_at":   time.Now(),
		})
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrDeviceNotFound
	}

	return &updated, nil
}

func (r *deviceRepository) MarkHubDevicesOffline(ctx context.Context, hubID uuid.UUID) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.Device{}).
		Where("hub_id = ? AND status <> ?", hubID, model.DeviceStatusOffline).
		Updates(map[string]any{
			"status":     model.DeviceStatusOffline,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
