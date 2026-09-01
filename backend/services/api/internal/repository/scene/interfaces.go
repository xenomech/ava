package scene

import (
	"context"

	"ava/api/internal/model"
	"ava/api/pkg/serrors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrSceneNotFound = serrors.NewCoded("scene_not_found", "scene not found")
	ErrNameTaken     = serrors.NewCoded("scene_name_taken", "this room already has a scene with that name")
)

type Repository interface {
	ListByRoom(ctx context.Context, tenantID, roomID uuid.UUID) ([]*model.Scene, error)
	GetByID(ctx context.Context, tenantID, roomID, sceneID uuid.UUID) (*model.Scene, error)
	Create(ctx context.Context, scene *model.Scene) error
	Delete(ctx context.Context, tenantID, roomID, sceneID uuid.UUID) error
	NextPosition(ctx context.Context, tenantID, roomID uuid.UUID) (int, error)
	NameExists(ctx context.Context, tenantID, roomID uuid.UUID, name string) (bool, error)
	DeviceIDsInRoom(ctx context.Context, tenantID, roomID uuid.UUID) ([]uuid.UUID, error)
}

type sceneRepository struct {
	db *gorm.DB
}

func NewRepository(database *gorm.DB) Repository {
	return &sceneRepository{db: database}
}
