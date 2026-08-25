package scene

import (
	"context"
	"encoding/json"
	"strings"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	roomrepo "ava/api/internal/repository/room"
	scenerepo "ava/api/internal/repository/scene"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"
	"ava/pkg/wire"

	"github.com/google/uuid"
)

var (
	ErrSceneNotFound = scenerepo.ErrSceneNotFound
	ErrNameTaken     = scenerepo.ErrNameTaken
	ErrRoomNotFound  = roomrepo.ErrRoomNotFound
	ErrNameRequired  = serrors.NewCoded("scene_name_required", "a scene needs a name")
	ErrNothingToSave = serrors.NewCoded("scene_empty", "a scene needs at least one device in this room")
)

func (s *sceneService) ListByRoom(
	ctx context.Context,
	tenantID, roomID uuid.UUID,
) ([]*dto.SceneResponse, error) {
	scenes, err := s.sceneRepo.ListByRoom(ctx, tenantID, roomID)
	if err != nil {
		logger.Error("scene.ListByRoom", logger.Err(err))

		return nil, err
	}

	out := make([]*dto.SceneResponse, 0, len(scenes))
	for _, found := range scenes {
		out = append(out, toSceneResponse(found))
	}

	return out, nil
}

// Create stores what the caller says the room is doing right now.
//
// The snapshot is taken on the client rather than read back out of the
// database, and deliberately: the person is saving the room they can see, which
// includes the brightness they finished dragging a moment ago and the hub has
// not confirmed yet. Reading the stored state instead would quietly save a
// slightly older room than the one on screen.
//
// What the server does insist on is that every device named is a device of
// this tenant, in this room. Capability checks are left to apply — that is the
// path with the authority, and a scene saved today should not become
// unsaveable because a bulb is unplugged.
func (s *sceneService) Create(
	ctx context.Context,
	tenantID, roomID uuid.UUID,
	req *dto.CreateSceneRequest,
) (*dto.SceneResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrNameRequired
	}

	if _, err := s.roomRepo.GetByID(ctx, tenantID, roomID); err != nil {
		if serrors.Is(err, roomrepo.ErrRoomNotFound) {
			return nil, ErrRoomNotFound
		}

		return nil, err
	}

	taken, err := s.sceneRepo.NameExists(ctx, tenantID, roomID, name)
	if err != nil {
		logger.Error("scene.Create", logger.Err(err))

		return nil, err
	}

	if taken {
		return nil, ErrNameTaken
	}

	inRoom, err := s.sceneRepo.DeviceIDsInRoom(ctx, tenantID, roomID)
	if err != nil {
		logger.Error("scene.Create", logger.Err(err))

		return nil, err
	}

	targets, err := plan(req.Targets, inRoom)
	if err != nil {
		return nil, err
	}

	if len(targets) == 0 {
		return nil, ErrNothingToSave
	}

	position, err := s.sceneRepo.NextPosition(ctx, tenantID, roomID)
	if err != nil {
		return nil, err
	}

	created := model.NewScene(tenantID, roomID, name, position)
	created.Targets = targets

	if err := s.sceneRepo.Create(ctx, created); err != nil {
		logger.Error("scene.Create", logger.Err(err))

		return nil, err
	}

	return toSceneResponse(created), nil
}

func (s *sceneService) Delete(ctx context.Context, tenantID, roomID, sceneID uuid.UUID) error {
	if err := s.sceneRepo.Delete(ctx, tenantID, roomID, sceneID); err != nil {
		if serrors.Is(err, scenerepo.ErrSceneNotFound) {
			return ErrSceneNotFound
		}

		logger.Error("scene.Delete", logger.Err(err))

		return err
	}

	return nil
}

// plan turns the requested targets into rows, dropping anything that is not a
// device of this room and keeping only the last value given for a trait.
func plan(wanted []dto.SceneTargetRequest, inRoom []uuid.UUID) ([]model.SceneTarget, error) {
	allowed := make(map[uuid.UUID]struct{}, len(inRoom))
	for _, id := range inRoom {
		allowed[id] = struct{}{}
	}

	type slot struct {
		device uuid.UUID
		trait  wire.Trait
	}

	at := make(map[slot]int, len(wanted))
	targets := make([]model.SceneTarget, 0, len(wanted))

	for index := range wanted {
		target := &wanted[index]

		if _, ok := allowed[target.DeviceID]; !ok {
			continue
		}

		if strings.TrimSpace(string(target.Trait)) == "" || !target.Value.IsSet() {
			continue
		}

		value, err := json.Marshal(target.Value)
		if err != nil {
			return nil, err
		}

		row := model.SceneTarget{
			DeviceID: target.DeviceID,
			Trait:    string(target.Trait),
			Value:    value,
		}

		key := slot{device: target.DeviceID, trait: target.Trait}
		if seen, duplicate := at[key]; duplicate {
			targets[seen] = row

			continue
		}

		at[key] = len(targets)
		targets = append(targets, row)
	}

	return targets, nil
}

func toSceneResponse(scene *model.Scene) *dto.SceneResponse {
	targets := make([]dto.SceneTargetResponse, 0, len(scene.Targets))

	for at := range scene.Targets {
		target := &scene.Targets[at]

		var value wire.Value
		if err := json.Unmarshal(target.Value, &value); err != nil {
			logger.Warn("SCENE_VALUE_UNREADABLE",
				logger.Any("scene.ID", scene.ID),
				logger.String("trait", target.Trait),
				logger.Err(err),
			)

			continue
		}

		targets = append(targets, dto.SceneTargetResponse{
			DeviceID: target.DeviceID,
			Trait:    wire.Trait(target.Trait),
			Value:    value,
		})
	}

	return &dto.SceneResponse{
		ID:       scene.ID,
		RoomID:   scene.RoomID,
		Name:     scene.Name,
		Position: scene.Position,
		Targets:  targets,
	}
}
