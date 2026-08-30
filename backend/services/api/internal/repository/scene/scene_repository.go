package scene

import (
	"context"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *sceneRepository) ListByRoom(ctx context.Context, tenantID, roomID uuid.UUID) ([]*model.Scene, error) {
	var scenes []*model.Scene

	err := r.db.WithContext(ctx).
		Preload("Targets").
		Where("tenant_id = ? AND room_id = ?", tenantID, roomID).
		Order("position ASC, created_at ASC").
		Find(&scenes).Error

	return scenes, err
}

// Create writes the scene and its targets together, so a scene never exists with half of the room in it.
func (r *sceneRepository) Create(ctx context.Context, scene *model.Scene) error {
	return r.db.WithContext(ctx).Create(scene).Error
}

func (r *sceneRepository) Delete(ctx context.Context, tenantID, roomID, sceneID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		result := dbTx.
			Where("tenant_id = ? AND room_id = ? AND id = ?", tenantID, roomID, sceneID).
			Delete(&model.Scene{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrSceneNotFound
		}

		// The cascade only fires on a hard delete, so a soft delete must sweep the targets by hand.
		return dbTx.Where("scene_id = ?", sceneID).Delete(&model.SceneTarget{}).Error
	})
}

func (r *sceneRepository) NextPosition(ctx context.Context, tenantID, roomID uuid.UUID) (int, error) {
	var highest *int

	err := r.db.WithContext(ctx).
		Model(&model.Scene{}).
		Where("tenant_id = ? AND room_id = ?", tenantID, roomID).
		Select("max(position)").
		Scan(&highest).Error
	if err != nil {
		return 0, err
	}

	if highest == nil {
		return 0, nil
	}

	return *highest + 1, nil
}

// NameExists asks under the default scope rather than relying on an index, so a deleted name is free again.
func (r *sceneRepository) NameExists(
	ctx context.Context,
	tenantID, roomID uuid.UUID,
	name string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.Scene{}).
		Where("tenant_id = ? AND room_id = ? AND lower(name) = lower(?)", tenantID, roomID, name).
		Count(&count).Error

	return count > 0, err
}

func (r *sceneRepository) DeviceIDsInRoom(
	ctx context.Context,
	tenantID, roomID uuid.UUID,
) ([]uuid.UUID, error) {
	var ids []uuid.UUID

	err := r.db.WithContext(ctx).
		Model(&model.Device{}).
		Where("tenant_id = ? AND room_id = ?", tenantID, roomID).
		Pluck("id", &ids).Error

	return ids, err
}
