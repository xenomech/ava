package scene

import (
	"context"

	"ava/api/internal/dto"
	roomrepo "ava/api/internal/repository/room"
	scenerepo "ava/api/internal/repository/scene"

	"github.com/google/uuid"
)

// Applier is the slice of the device service a scene needs: the batch write it replays through.
type Applier interface {
	Apply(ctx context.Context, tenantID uuid.UUID, req *dto.ApplyRequest) (*dto.ApplyResponse, error)
}

type Service interface {
	ListByRoom(ctx context.Context, tenantID, roomID uuid.UUID) ([]*dto.SceneResponse, error)
	Create(ctx context.Context, tenantID, roomID uuid.UUID, req *dto.CreateSceneRequest) (*dto.SceneResponse, error)
	Delete(ctx context.Context, tenantID, roomID, sceneID uuid.UUID) error
	Apply(ctx context.Context, tenantID, roomID, sceneID uuid.UUID) (*dto.ApplyResponse, error)
}

type sceneService struct {
	sceneRepo scenerepo.Repository
	roomRepo  roomrepo.Repository
	applier   Applier
}

func NewService(
	sceneRepo scenerepo.Repository,
	roomRepo roomrepo.Repository,
	applier Applier,
) Service {
	return &sceneService{sceneRepo: sceneRepo, roomRepo: roomRepo, applier: applier}
}
