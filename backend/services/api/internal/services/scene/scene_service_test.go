package scene

import (
	"context"
	"encoding/json"
	"testing"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	roomrepo "ava/api/internal/repository/room"
	scenerepo "ava/api/internal/repository/scene"
	"ava/api/pkg/serrors"
	"ava/pkg/wire"

	"github.com/google/uuid"
)

type fakeScenes struct {
	scenes  []*model.Scene
	inRoom  []uuid.UUID
	taken   bool
	next    int
	created *model.Scene
	deleted uuid.UUID
	fail    error
}

func (f *fakeScenes) ListByRoom(_ context.Context, _, _ uuid.UUID) ([]*model.Scene, error) {
	return f.scenes, f.fail
}

func (f *fakeScenes) GetByID(_ context.Context, _, _, _ uuid.UUID) (*model.Scene, error) {
	if f.fail != nil {
		return nil, f.fail
	}

	if len(f.scenes) == 0 {
		return nil, scenerepo.ErrSceneNotFound
	}

	return f.scenes[0], nil
}

func (f *fakeScenes) Create(_ context.Context, scene *model.Scene) error {
	if f.fail != nil {
		return f.fail
	}

	f.created = scene

	return nil
}

func (f *fakeScenes) Delete(_ context.Context, _, _, sceneID uuid.UUID) error {
	if f.fail != nil {
		return f.fail
	}

	f.deleted = sceneID

	return nil
}

func (f *fakeScenes) NextPosition(_ context.Context, _, _ uuid.UUID) (int, error) {
	return f.next, nil
}

func (f *fakeScenes) NameExists(_ context.Context, _, _ uuid.UUID, _ string) (bool, error) {
	return f.taken, nil
}

func (f *fakeScenes) DeviceIDsInRoom(_ context.Context, _, _ uuid.UUID) ([]uuid.UUID, error) {
	return f.inRoom, nil
}

// fakeApplier records what a scene handed to the batch write.
type fakeApplier struct {
	got  []dto.ApplyTargetRequest
	fail error
}

func (f *fakeApplier) Apply(
	_ context.Context,
	_ uuid.UUID,
	req *dto.ApplyRequest,
) (*dto.ApplyResponse, error) {
	if f.fail != nil {
		return nil, f.fail
	}

	f.got = req.Targets

	return &dto.ApplyResponse{}, nil
}

type fakeRooms struct {
	missing bool
}

func (f *fakeRooms) ListByTenant(_ context.Context, _ uuid.UUID) ([]*model.Room, error) {
	return nil, nil
}

func (f *fakeRooms) GetByID(_ context.Context, _, roomID uuid.UUID) (*model.Room, error) {
	if f.missing {
		return nil, roomrepo.ErrRoomNotFound
	}

	return &model.Room{Name: "Living Room"}, nil
}

func (f *fakeRooms) Create(_ context.Context, _ *model.Room) error { return nil }

func (f *fakeRooms) Update(_ context.Context, _, _ uuid.UUID, _ map[string]any) (*model.Room, error) {
	return &model.Room{}, nil
}

func (f *fakeRooms) Delete(_ context.Context, _, _ uuid.UUID) error { return nil }

func (f *fakeRooms) NextPosition(_ context.Context, _ uuid.UUID) (int, error) { return 0, nil }

func request(name string, targets ...dto.SceneTargetRequest) *dto.CreateSceneRequest {
	return &dto.CreateSceneRequest{Name: name, Targets: targets}
}

func target(device uuid.UUID, trait wire.Trait, value wire.Value) dto.SceneTargetRequest {
	return dto.SceneTargetRequest{DeviceID: device, Trait: trait, Value: value}
}

func save(t *testing.T, repo *fakeScenes, req *dto.CreateSceneRequest) error {
	t.Helper()

	_, err := NewService(repo, &fakeRooms{}, &fakeApplier{}).
		Create(context.Background(), uuid.New(), uuid.New(), req)

	return err
}

func TestASceneOnlyKeepsDevicesThatAreInTheRoom(t *testing.T) {
	lamp, elsewhere := uuid.New(), uuid.New()
	repo := &fakeScenes{inRoom: []uuid.UUID{lamp}}

	err := save(t, repo, request("Evening",
		target(lamp, wire.TraitPower, wire.Bool(true)),
		target(elsewhere, wire.TraitPower, wire.Bool(true)),
	))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if len(repo.created.Targets) != 1 {
		t.Fatalf("kept %d targets, want only the one in the room", len(repo.created.Targets))
	}

	if repo.created.Targets[0].DeviceID != lamp {
		t.Errorf("kept %v, want the lamp", repo.created.Targets[0].DeviceID)
	}
}

func TestASceneWithNothingLeftInTheRoomIsRefused(t *testing.T) {
	repo := &fakeScenes{inRoom: []uuid.UUID{uuid.New()}}

	err := save(t, repo, request("Evening", target(uuid.New(), wire.TraitPower, wire.Bool(true))))
	if !serrors.Is(err, ErrNothingToSave) {
		t.Errorf("got %v, want ErrNothingToSave", err)
	}

	if repo.created != nil {
		t.Error("an empty scene reached the repository")
	}
}

func TestTheLastValueForATraitWins(t *testing.T) {
	lamp := uuid.New()
	repo := &fakeScenes{inRoom: []uuid.UUID{lamp}}

	err := save(t, repo, request("Evening",
		target(lamp, wire.TraitBrightness, wire.Number(20)),
		target(lamp, wire.TraitBrightness, wire.Number(80)),
	))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if len(repo.created.Targets) != 1 {
		t.Fatalf("stored %d targets, want 1", len(repo.created.Targets))
	}

	var stored wire.Value
	if err := json.Unmarshal(repo.created.Targets[0].Value, &stored); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if number, _ := stored.Number(); number != 80 {
		t.Errorf("brightness = %v, want the later 80", number)
	}
}

func TestAnUnsetValueIsNotWorthStoring(t *testing.T) {
	lamp := uuid.New()
	repo := &fakeScenes{inRoom: []uuid.UUID{lamp}}

	err := save(t, repo, request("Evening",
		target(lamp, wire.TraitPower, wire.Bool(false)),
		target(lamp, wire.TraitColor, wire.Value{}),
	))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if len(repo.created.Targets) != 1 {
		t.Errorf("stored %d targets, want only the one that has a value", len(repo.created.Targets))
	}
}

func TestASceneNameIsTrimmedAndCannotBeBlankOrTaken(t *testing.T) {
	lamp := uuid.New()
	power := target(lamp, wire.TraitPower, wire.Bool(true))

	repo := &fakeScenes{inRoom: []uuid.UUID{lamp}}
	if err := save(t, repo, request("  Evening  ", power)); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if repo.created.Name != "Evening" {
		t.Errorf("name = %q, want it trimmed", repo.created.Name)
	}

	blank := &fakeScenes{inRoom: []uuid.UUID{lamp}}
	if err := save(t, blank, request("   ", power)); !serrors.Is(err, ErrNameRequired) {
		t.Errorf("a blank name gave %v", err)
	}

	twice := &fakeScenes{inRoom: []uuid.UUID{lamp}, taken: true}
	if err := save(t, twice, request("Evening", power)); !serrors.Is(err, ErrNameTaken) {
		t.Errorf("a duplicate name gave %v", err)
	}
}

func TestANewSceneGoesToTheEndOfTheRow(t *testing.T) {
	lamp := uuid.New()
	repo := &fakeScenes{inRoom: []uuid.UUID{lamp}, next: 3}

	if err := save(t, repo, request("Evening", target(lamp, wire.TraitPower, wire.Bool(true)))); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if repo.created.Position != 3 {
		t.Errorf("position = %d, want 3", repo.created.Position)
	}
}

func TestSavingIntoARoomThatIsGoneIsNotFound(t *testing.T) {
	lamp := uuid.New()
	service := NewService(&fakeScenes{inRoom: []uuid.UUID{lamp}}, &fakeRooms{missing: true}, &fakeApplier{})

	_, err := service.Create(context.Background(), uuid.New(), uuid.New(),
		request("Evening", target(lamp, wire.TraitPower, wire.Bool(true))))
	if !serrors.Is(err, ErrRoomNotFound) {
		t.Errorf("got %v, want ErrRoomNotFound", err)
	}
}

func TestDeletingASceneThatIsNotThereIsReportedAsNotFound(t *testing.T) {
	service := NewService(&fakeScenes{fail: scenerepo.ErrSceneNotFound}, &fakeRooms{}, &fakeApplier{})

	err := service.Delete(context.Background(), uuid.New(), uuid.New(), uuid.New())
	if !serrors.Is(err, ErrSceneNotFound) {
		t.Errorf("got %v", err)
	}
}

func TestATargetWhoseValueCannotBeReadIsLeftOutOfTheResponse(t *testing.T) {
	lamp := uuid.New()
	repo := &fakeScenes{scenes: []*model.Scene{{
		Name: "Evening",
		Targets: []model.SceneTarget{
			{DeviceID: lamp, Trait: string(wire.TraitPower), Value: json.RawMessage(`true`)},
			{DeviceID: lamp, Trait: string(wire.TraitColor), Value: json.RawMessage(`{"broken":`)},
		},
	}}}

	scenes, err := NewService(repo, &fakeRooms{}, &fakeApplier{}).ListByRoom(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("ListByRoom: %v", err)
	}

	if len(scenes) != 1 || len(scenes[0].Targets) != 1 {
		t.Fatalf("got %+v, want the readable target only", scenes)
	}

	if scenes[0].Targets[0].Trait != wire.TraitPower {
		t.Errorf("kept %q", scenes[0].Targets[0].Trait)
	}
}

func sceneWith(targets ...model.SceneTarget) *model.Scene {
	scene := model.NewScene(uuid.New(), uuid.New(), "Evening", 0)
	scene.Targets = targets

	return scene
}

func savedTarget(device uuid.UUID, trait, raw string) model.SceneTarget {
	return model.SceneTarget{DeviceID: device, Trait: trait, Value: []byte(raw)}
}

func TestApplyingASceneSendsItsTargetsToTheBatchWrite(t *testing.T) {
	lamp := uuid.New()
	repo := &fakeScenes{scenes: []*model.Scene{sceneWith(
		savedTarget(lamp, "brightness", "40"),
		savedTarget(lamp, "power", "true"),
	)}}
	applier := &fakeApplier{}

	if _, err := NewService(repo, &fakeRooms{}, applier).
		Apply(context.Background(), uuid.New(), uuid.New(), uuid.New()); err != nil {
		t.Fatalf("apply failed: %v", err)
	}

	if len(applier.got) != 2 {
		t.Fatalf("sent %d targets, wanted 2", len(applier.got))
	}

	if applier.got[0].Trait != "brightness" || applier.got[1].Trait != "power" {
		t.Fatalf("targets arrived as %v", applier.got)
	}
}

func TestApplyingAMissingSceneIsNotFound(t *testing.T) {
	service := NewService(&fakeScenes{}, &fakeRooms{}, &fakeApplier{})

	_, err := service.Apply(context.Background(), uuid.New(), uuid.New(), uuid.New())
	if !serrors.Is(err, ErrSceneNotFound) {
		t.Fatalf("got %v, wanted ErrSceneNotFound", err)
	}
}

func TestApplyingAnEmptySceneSendsNothing(t *testing.T) {
	applier := &fakeApplier{}
	service := NewService(&fakeScenes{scenes: []*model.Scene{sceneWith()}}, &fakeRooms{}, applier)

	_, err := service.Apply(context.Background(), uuid.New(), uuid.New(), uuid.New())
	if !serrors.Is(err, ErrNothingToApply) {
		t.Fatalf("got %v, wanted ErrNothingToApply", err)
	}

	if applier.got != nil {
		t.Fatal("an empty scene still reached the batch write")
	}
}

// A target written by an older version must not take the whole scene down with it.
func TestASceneWithOneUnreadableTargetStillAppliesTheRest(t *testing.T) {
	lamp := uuid.New()
	repo := &fakeScenes{scenes: []*model.Scene{sceneWith(
		savedTarget(lamp, "brightness", "not json"),
		savedTarget(lamp, "power", "true"),
	)}}
	applier := &fakeApplier{}

	if _, err := NewService(repo, &fakeRooms{}, applier).
		Apply(context.Background(), uuid.New(), uuid.New(), uuid.New()); err != nil {
		t.Fatalf("apply failed: %v", err)
	}

	if len(applier.got) != 1 || applier.got[0].Trait != "power" {
		t.Fatalf("expected only the readable target, got %v", applier.got)
	}
}
