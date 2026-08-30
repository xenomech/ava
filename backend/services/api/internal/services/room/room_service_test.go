package room

import (
	"context"
	"testing"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	roomrepo "ava/api/internal/repository/room"
	"ava/api/pkg/serrors"

	"github.com/google/uuid"
)

type fakeRepo struct {
	rooms    []*model.Room
	next     int
	created  *model.Room
	updated  map[string]any
	deleted  uuid.UUID
	failWith error
}

func (f *fakeRepo) ListByTenant(_ context.Context, _ uuid.UUID) ([]*model.Room, error) {
	return f.rooms, f.failWith
}

func (f *fakeRepo) GetByID(_ context.Context, _, roomID uuid.UUID) (*model.Room, error) {
	for _, room := range f.rooms {
		if room.ID == roomID {
			return room, nil
		}
	}

	return nil, roomrepo.ErrRoomNotFound
}

func (f *fakeRepo) Create(_ context.Context, room *model.Room) error {
	if f.failWith != nil {
		return f.failWith
	}

	f.created = room

	return nil
}

func (f *fakeRepo) Update(_ context.Context, _, _ uuid.UUID, fields map[string]any) (*model.Room, error) {
	if f.failWith != nil {
		return nil, f.failWith
	}

	f.updated = fields

	return &model.Room{Name: "after"}, nil
}

func (f *fakeRepo) Delete(_ context.Context, _, roomID uuid.UUID) error {
	if f.failWith != nil {
		return f.failWith
	}

	f.deleted = roomID

	return nil
}

func (f *fakeRepo) NextPosition(_ context.Context, _ uuid.UUID) (int, error) {
	return f.next, nil
}

func text(v string) *string { return &v }

func TestANewRoomGoesToTheEndOfTheList(t *testing.T) {
	repo := &fakeRepo{next: 4}
	service := NewService(repo)

	if _, err := service.Create(context.Background(), uuid.New(), &dto.CreateRoomRequest{Name: "Kitchen"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if repo.created.Position != 4 {
		t.Errorf("position = %d, want 4", repo.created.Position)
	}
}

func TestARoomNameIsTrimmedAndCannotBeBlank(t *testing.T) {
	repo := &fakeRepo{}
	service := NewService(repo)

	if _, err := service.Create(context.Background(), uuid.New(), &dto.CreateRoomRequest{Name: "  Loft  "}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if repo.created.Name != "Loft" {
		t.Errorf("name = %q, want it trimmed", repo.created.Name)
	}

	_, err := service.Create(context.Background(), uuid.New(), &dto.CreateRoomRequest{Name: "   "})
	if !serrors.Is(err, ErrNameRequired) {
		t.Errorf("a blank name gave %v", err)
	}
}

func TestADuplicateNameSurfacesAsItsOwnError(t *testing.T) {
	service := NewService(&fakeRepo{failWith: roomrepo.ErrNameTaken})

	_, err := service.Create(context.Background(), uuid.New(), &dto.CreateRoomRequest{Name: "Kitchen"})
	if !serrors.Is(err, ErrNameTaken) {
		t.Errorf("got %v", err)
	}
}

func TestUpdateOnlyTouchesTheFieldsThatWereSent(t *testing.T) {
	position := 2

	cases := map[string]struct {
		request dto.UpdateRoomRequest
		fields  []string
	}{
		"rename only":   {dto.UpdateRoomRequest{Name: text("Lounge")}, []string{"name"}},
		"reorder only":  {dto.UpdateRoomRequest{Position: &position}, []string{"position"}},
		"both together": {dto.UpdateRoomRequest{Name: text("Lounge"), Position: &position}, []string{"name", "position"}},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			repo := &fakeRepo{}
			service := NewService(repo)

			if _, err := service.Update(context.Background(), uuid.New(), uuid.New(), &tc.request); err != nil {
				t.Fatalf("Update: %v", err)
			}

			for _, field := range tc.fields {
				if _, ok := repo.updated[field]; !ok {
					t.Errorf("%s was not written", field)
				}
			}

			if len(repo.updated) != len(tc.fields) {
				t.Errorf("wrote %v, want only %v", repo.updated, tc.fields)
			}
		})
	}
}

func TestRenamingToBlankIsRejectedBeforeItReachesTheDatabase(t *testing.T) {
	repo := &fakeRepo{}
	service := NewService(repo)

	_, err := service.Update(context.Background(), uuid.New(), uuid.New(), &dto.UpdateRoomRequest{Name: text("  ")})
	if !serrors.Is(err, ErrNameRequired) {
		t.Errorf("got %v", err)
	}

	if repo.updated != nil {
		t.Error("a blank rename reached the repository")
	}
}

func TestAMissingRoomIsReportedAsNotFound(t *testing.T) {
	service := NewService(&fakeRepo{failWith: roomrepo.ErrRoomNotFound})

	if err := service.Delete(context.Background(), uuid.New(), uuid.New()); !serrors.Is(err, ErrRoomNotFound) {
		t.Errorf("Delete gave %v", err)
	}

	if _, err := service.Update(context.Background(), uuid.New(), uuid.New(),
		&dto.UpdateRoomRequest{Name: text("Lounge")}); !serrors.Is(err, ErrRoomNotFound) {
		t.Errorf("Update gave %v", err)
	}
}

func TestRoomsAreReturnedInTheOrderTheRepositoryGivesThem(t *testing.T) {
	repo := &fakeRepo{rooms: []*model.Room{
		{Name: "Hallway", Position: 0},
		{Name: "Kitchen", Position: 1},
	}}

	rooms, err := NewService(repo).ListByTenant(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("ListByTenant: %v", err)
	}

	if len(rooms) != 2 || rooms[0].Name != "Hallway" || rooms[1].Name != "Kitchen" {
		t.Errorf("got %+v", rooms)
	}
}
