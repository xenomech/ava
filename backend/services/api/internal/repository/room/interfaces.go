package room

import (
	"context"

	"ava/api/internal/model"
	"ava/api/pkg/serrors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRoomNotFound = serrors.NewCoded("room_not_found", "room not found")
	ErrNameTaken    = serrors.NewCoded("room_name_taken", "a room with that name already exists")
)

type Repository interface {
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*model.Room, error)
	GetByID(ctx context.Context, tenantID, roomID uuid.UUID) (*model.Room, error)
	Create(ctx context.Context, room *model.Room) error
	Update(ctx context.Context, tenantID, roomID uuid.UUID, fields map[string]any) (*model.Room, error)
	Delete(ctx context.Context, tenantID, roomID uuid.UUID) error
	NextPosition(ctx context.Context, tenantID uuid.UUID) (int, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRepository(database *gorm.DB) Repository {
	return &roomRepository{db: database}
}
