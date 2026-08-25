package scene

import (
	"context"

	"ava/api/internal/dto"
	roomrepo "ava/api/internal/repository/room"
	scenerepo "ava/api/internal/repository/scene"

	"github.com/google/uuid"
)

type Service interface {
	ListByRoom(ctx context.Context, tenantID, roomID uuid.UUID) ([]*dto.SceneResponse, error)
	Create(ctx context.Context, tenantID, roomID uuid.UUID, req *dto.CreateSceneRequest) (*dto.SceneResponse, error)
	Delete(ctx context.Context, tenantID, roomID, sceneID uuid.UUID) error
}

type sceneService struct {
	sceneRepo scenerepo.Repository
	roomRepo  roomrepo.Repository
}

func NewService(sceneRepo scenerepo.Repository, roomRepo roomrepo.Repository) Service {
	return &sceneService{sceneRepo: sceneRepo, roomRepo: roomRepo}
}
