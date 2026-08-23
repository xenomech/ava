package room_test

import (
	"context"
	"os"
	"testing"

	"ava/api/internal/db"
	"ava/api/internal/model"
	roomrepo "ava/api/internal/repository/room"
	"ava/api/pkg/serrors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func connect(t *testing.T) *gorm.DB {
	t.Helper()

	host := os.Getenv("AVA_TEST_DB_HOST")
	if host == "" {
		t.Skip("set AVA_TEST_DB_HOST, AVA_TEST_DB_PORT, AVA_TEST_DB_USER, AVA_TEST_DB_PASSWORD and AVA_TEST_DB_NAME to run")
	}

	database, err := db.Connect(&db.PostgresConfig{
		Host:     host,
		Port:     os.Getenv("AVA_TEST_DB_PORT"),
		User:     os.Getenv("AVA_TEST_DB_USER"),
		Password: os.Getenv("AVA_TEST_DB_PASSWORD"),
		Database: os.Getenv("AVA_TEST_DB_NAME"),
		SSLMode:  "disable",
	})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}

	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	return database
}

func tenant(t *testing.T, database *gorm.DB) uuid.UUID {
	t.Helper()

	created := &model.Tenant{Name: "Room Repo Test", Slug: "room-repo-" + uuid.NewString()[:8]}
	if err := database.Create(created).Error; err != nil {
		t.Fatalf("create tenant: %v", err)
	}

	t.Cleanup(func() { database.Unscoped().Delete(&model.Tenant{}, "id = ?", created.ID) })

	return created.ID
}

func TestDeletingARoomLeavesItsDevicesBehindUnassigned(t *testing.T) {
	database := connect(t)
	tenantID := tenant(t, database)
	repo := roomrepo.NewRepository(database)
	ctx := context.Background()

	room := model.NewRoom(tenantID, "Living Room", 0)
	if err := repo.Create(ctx, room); err != nil {
		t.Fatalf("Create: %v", err)
	}

	hub := &model.Hub{TenantID: tenantID, Name: "test hub", Status: model.HubStatusActive}
	if err := database.Create(hub).Error; err != nil {
		t.Fatalf("create hub: %v", err)
	}

	device := model.NewDevice(tenantID, hub.ID, &model.Reported{
		ExternalID: "lamp-" + uuid.NewString()[:8],
		Name:       "Lamp",
		Kind:       "bulb",
		Status:     model.DeviceStatusOnline,
	})
	device.RoomID = &room.ID

	if err := database.Create(device).Error; err != nil {
		t.Fatalf("create device: %v", err)
	}

	if err := repo.Delete(ctx, tenantID, room.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	var after model.Device
	if err := database.First(&after, "id = ?", device.ID).Error; err != nil {
		t.Fatalf("the device was deleted along with its room: %v", err)
	}

	if after.RoomID != nil {
		t.Errorf("room_id = %v, want nil — a soft deleted room must not keep devices pointing at it", after.RoomID)
	}
}

func TestTwoRoomsInOneTenantCannotShareAName(t *testing.T) {
	database := connect(t)
	tenantID := tenant(t, database)
	repo := roomrepo.NewRepository(database)
	ctx := context.Background()

	if err := repo.Create(ctx, model.NewRoom(tenantID, "Kitchen", 0)); err != nil {
		t.Fatalf("Create: %v", err)
	}

	err := repo.Create(ctx, model.NewRoom(tenantID, "Kitchen", 1))
	if !serrors.Is(err, roomrepo.ErrNameTaken) {
		t.Errorf("got %v, want ErrNameTaken", err)
	}
}

func TestRoomsComeBackInPositionOrder(t *testing.T) {
	database := connect(t)
	tenantID := tenant(t, database)
	repo := roomrepo.NewRepository(database)
	ctx := context.Background()

	for name, position := range map[string]int{"Attic": 2, "Hallway": 0, "Kitchen": 1} {
		if err := repo.Create(ctx, model.NewRoom(tenantID, name, position)); err != nil {
			t.Fatalf("Create %s: %v", name, err)
		}
	}

	rooms, err := repo.ListByTenant(ctx, tenantID)
	if err != nil {
		t.Fatalf("ListByTenant: %v", err)
	}

	got := make([]string, 0, len(rooms))
	for _, room := range rooms {
		got = append(got, room.Name)
	}

	if len(got) != 3 || got[0] != "Hallway" || got[1] != "Kitchen" || got[2] != "Attic" {
		t.Errorf("order = %v", got)
	}
}

func TestAnotherTenantsRoomIsInvisible(t *testing.T) {
	database := connect(t)
	mine := tenant(t, database)
	theirs := tenant(t, database)
	repo := roomrepo.NewRepository(database)
	ctx := context.Background()

	room := model.NewRoom(theirs, "Their Kitchen", 0)
	if err := repo.Create(ctx, room); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if _, err := repo.GetByID(ctx, mine, room.ID); !serrors.Is(err, roomrepo.ErrRoomNotFound) {
		t.Errorf("GetByID across tenants gave %v", err)
	}

	if err := repo.Delete(ctx, mine, room.ID); !serrors.Is(err, roomrepo.ErrRoomNotFound) {
		t.Errorf("Delete across tenants gave %v", err)
	}
}
