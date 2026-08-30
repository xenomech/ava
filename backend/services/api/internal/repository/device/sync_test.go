package device_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"ava/api/internal/db"
	"ava/api/internal/model"
	devicerepo "ava/api/internal/repository/device"

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

func scaffold(t *testing.T, database *gorm.DB) (tenantID, hubID uuid.UUID) {
	t.Helper()

	tenant := &model.Tenant{Name: "Sync Test", Slug: "sync-" + uuid.NewString()[:8]}
	if err := database.Create(tenant).Error; err != nil {
		t.Fatalf("create tenant: %v", err)
	}

	t.Cleanup(func() { database.Unscoped().Delete(&model.Tenant{}, "id = ?", tenant.ID) })

	hub := &model.Hub{TenantID: tenant.ID, Name: "test hub", Status: model.HubStatusActive}
	if err := database.Create(hub).Error; err != nil {
		t.Fatalf("create hub: %v", err)
	}

	return tenant.ID, hub.ID
}

func report(tenantID, hubID uuid.UUID, externalID, state string) []*model.Device {
	return []*model.Device{
		model.NewDevice(tenantID, hubID, &model.Reported{
			ExternalID: externalID,
			Name:       "Lamp",
			Kind:       "bulb",
			Status:     model.DeviceStatusOnline,
			State:      json.RawMessage(state),
		}),
	}
}

func stateOf(t *testing.T, database *gorm.DB, externalID string) map[string]any {
	t.Helper()

	var found model.Device
	if err := database.First(&found, "external_id = ?", externalID).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}

	var state map[string]any
	if err := json.Unmarshal(found.State, &state); err != nil {
		t.Fatalf("decode state: %v", err)
	}

	return state
}

// A sweep carries a reading seconds old, so writing it over stored state would undo a newer change.
func TestASweepDoesNotOverwriteWhatADeviceIsAlreadyDoing(t *testing.T) {
	database := connect(t)
	tenantID, hubID := scaffold(t, database)
	repo := devicerepo.NewRepository(database)
	ctx := context.Background()

	id := "lamp-" + uuid.NewString()[:8]

	if err := repo.SyncHubDevices(ctx, tenantID, hubID, report(tenantID, hubID, id, `{"power":true,"color_temp":2700}`)); err != nil {
		t.Fatalf("first sync: %v", err)
	}

	// Somebody sets a colour. This is the merging path, the way real state arrives.
	if _, err := repo.ApplyState(ctx, hubID, id, json.RawMessage(`{"color":"#ff3366"}`), []string{"color_temp"}); err != nil {
		t.Fatalf("apply state: %v", err)
	}

	// A sweep that began before that change now reports what it saw back then.
	if err := repo.SyncHubDevices(ctx, tenantID, hubID, report(tenantID, hubID, id, `{"power":true,"color_temp":2700}`)); err != nil {
		t.Fatalf("stale sync: %v", err)
	}

	state := stateOf(t, database, id)

	if state["color"] != "#ff3366" {
		t.Errorf("colour = %v, want it kept — a stale sweep reverted a newer change", state["color"])
	}

	if _, revived := state["color_temp"]; revived {
		t.Error("color_temp came back; the bulb is showing a colour and cannot hold both")
	}
}

// A device nobody has seen before has nothing to protect, so its first reading has to land.
func TestADeviceSeenForTheFirstTimeKeepsItsReading(t *testing.T) {
	database := connect(t)
	tenantID, hubID := scaffold(t, database)
	repo := devicerepo.NewRepository(database)

	id := "lamp-" + uuid.NewString()[:8]

	err := repo.SyncHubDevices(
		context.Background(),
		tenantID, hubID,
		report(tenantID, hubID, id, `{"power":true,"brightness":80,"color":"#4bd6ff"}`),
	)
	if err != nil {
		t.Fatalf("sync: %v", err)
	}

	state := stateOf(t, database, id)

	if state["color"] != "#4bd6ff" || state["brightness"] != float64(80) {
		t.Errorf("state = %v, want the reading it was discovered with", state)
	}
}
