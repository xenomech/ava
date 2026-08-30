package room

import (
	"context"

	"ava/api/internal/dto"
	roomrepo "ava/api/internal/repository/room"

	"github.com/google/uuid"
)

type Service interface {
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.RoomResponse, error)
	Create(ctx context.Context, tenantID uuid.UUID, req *dto.CreateRoomRequest) (*dto.RoomResponse, error)
	Update(ctx context.Context, tenantID, roomID uuid.UUID, req *dto.UpdateRoomRequest) (*dto.RoomResponse, error)
	Delete(ctx context.Context, tenantID, roomID uuid.UUID) error
}

type roomService struct {
	roomRepo roomrepo.Repository
}

func NewService(roomRepo roomrepo.Repository) Service {
	return &roomService{roomRepo: roomRepo}
}
