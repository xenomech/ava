package room

import (
	"context"
	"errors"
	"strings"
	"time"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *roomRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*model.Room, error) {
	var rooms []*model.Room

	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("position ASC, name ASC").
		Find(&rooms).Error

	return rooms, err
}

func (r *roomRepository) GetByID(ctx context.Context, tenantID, roomID uuid.UUID) (*model.Room, error) {
	var found model.Room

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, roomID).
		First(&found).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}

		return nil, err
	}

	return &found, nil
}

func (r *roomRepository) Create(ctx context.Context, room *model.Room) error {
	err := r.db.WithContext(ctx).Create(room).Error
	if err != nil && isDuplicateName(err) {
		return ErrNameTaken
	}

	return err
}

func (r *roomRepository) Update(
	ctx context.Context,
	tenantID, roomID uuid.UUID,
	fields map[string]any,
) (*model.Room, error) {
	if len(fields) == 0 {
		return r.GetByID(ctx, tenantID, roomID)
	}

	fields["updated_at"] = time.Now()

	var updated model.Room

	result := r.db.WithContext(ctx).Model(&updated).
		Clauses(clause.Returning{}).
		Where("tenant_id = ? AND id = ?", tenantID, roomID).
		Updates(fields)
	if result.Error != nil {
		if isDuplicateName(result.Error) {
			return nil, ErrNameTaken
		}

		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrRoomNotFound
	}

	return &updated, nil
}

func (r *roomRepository) Delete(ctx context.Context, tenantID, roomID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		err := dbTx.Model(&model.Device{}).
			Where("tenant_id = ? AND room_id = ?", tenantID, roomID).
			Updates(map[string]any{"room_id": nil, "updated_at": time.Now()}).Error
		if err != nil {
			return err
		}

		result := dbTx.Where("tenant_id = ? AND id = ?", tenantID, roomID).Delete(&model.Room{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrRoomNotFound
		}

		return nil
	})
}

func (r *roomRepository) NextPosition(ctx context.Context, tenantID uuid.UUID) (int, error) {
	var highest *int

	err := r.db.WithContext(ctx).
		Model(&model.Room{}).
		Where("tenant_id = ?", tenantID).
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

func isDuplicateName(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "idx_room_tenant_name")
}
