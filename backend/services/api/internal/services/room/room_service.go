package room

import (
	"context"
	"strings"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	roomrepo "ava/api/internal/repository/room"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"

	"github.com/google/uuid"
)

var (
	ErrRoomNotFound = roomrepo.ErrRoomNotFound
	ErrNameTaken    = roomrepo.ErrNameTaken
	ErrNameRequired = serrors.NewCoded("room_name_required", "a room needs a name")
)

func (s *roomService) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.RoomResponse, error) {
	rooms, err := s.roomRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		logger.Error("room.ListByTenant", logger.Err(err))

		return nil, err
	}

	return toRoomResponses(rooms), nil
}

func (s *roomService) Create(
	ctx context.Context,
	tenantID uuid.UUID,
	req *dto.CreateRoomRequest,
) (*dto.RoomResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrNameRequired
	}

	position, err := s.roomRepo.NextPosition(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	created := model.NewRoom(tenantID, name, position)

	if err := s.roomRepo.Create(ctx, created); err != nil {
		if serrors.Is(err, roomrepo.ErrNameTaken) {
			return nil, ErrNameTaken
		}

		logger.Error("room.Create", logger.Err(err))

		return nil, err
	}

	return toRoomResponse(created), nil
}

func (s *roomService) Update(
	ctx context.Context,
	tenantID, roomID uuid.UUID,
	req *dto.UpdateRoomRequest,
) (*dto.RoomResponse, error) {
	fields := make(map[string]any, 2)

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, ErrNameRequired
		}

		fields["name"] = name
	}

	if req.Position != nil {
		fields["position"] = *req.Position
	}

	updated, err := s.roomRepo.Update(ctx, tenantID, roomID, fields)
	if err != nil {
		if serrors.Is(err, roomrepo.ErrNameTaken) {
			return nil, ErrNameTaken
		}

		if serrors.Is(err, roomrepo.ErrRoomNotFound) {
			return nil, ErrRoomNotFound
		}

		logger.Error("room.Update", logger.Err(err))

		return nil, err
	}

	return toRoomResponse(updated), nil
}

func (s *roomService) Delete(ctx context.Context, tenantID, roomID uuid.UUID) error {
	if err := s.roomRepo.Delete(ctx, tenantID, roomID); err != nil {
		if serrors.Is(err, roomrepo.ErrRoomNotFound) {
			return ErrRoomNotFound
		}

		logger.Error("room.Delete", logger.Err(err))

		return err
	}

	return nil
}

func toRoomResponse(room *model.Room) *dto.RoomResponse {
	return &dto.RoomResponse{
		ID:       room.ID,
		Name:     room.Name,
		Position: room.Position,
	}
}

func toRoomResponses(rooms []*model.Room) []*dto.RoomResponse {
	out := make([]*dto.RoomResponse, 0, len(rooms))
	for _, room := range rooms {
		out = append(out, toRoomResponse(room))
	}

	return out
}
