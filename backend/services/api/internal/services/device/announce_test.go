package device

import (
	"testing"

	"ava/api/internal/dto"

	"github.com/google/uuid"
)

// A hub sweeps every thirty seconds whether or not anything happened, and this
// list replaces a client's whole view of that hub. Repeating an identical one
// is pure interference.
func TestAnUnchangedSweepIsNotBroadcast(t *testing.T) {
	service := &deviceService{}
	hub := uuid.New()

	list := []byte(`{"type":"device.list","devices":[{"id":"a"}]}`)

	if !service.listChanged(hub, list) {
		t.Fatal("the first list should always go out")
	}

	if service.listChanged(hub, list) {
		t.Error("an identical sweep was broadcast again")
	}

	if !service.listChanged(hub, []byte(`{"type":"device.list","devices":[{"id":"b"}]}`)) {
		t.Error("a changed list was suppressed — clients would never hear about it")
	}
}

func TestEachHubIsRememberedSeparately(t *testing.T) {
	service := &deviceService{}
	one, other := uuid.New(), uuid.New()

	list := []byte(`{"devices":[]}`)

	service.listChanged(one, list)

	if !service.listChanged(other, list) {
		t.Error("a second hub was silenced by the first hub's list")
	}
}

// Nothing above should stop a genuine announcement reaching the event service.
func TestAChangedListStillReachesClients(t *testing.T) {
	events := &recordingEvents{}
	service := &deviceService{events: events}
	hub := uuid.New()

	service.announceList(uuid.New(), hub, []*dto.DeviceResponse{{Name: "Lamp"}})
	service.announceList(uuid.New(), hub, []*dto.DeviceResponse{{Name: "Lamp"}})
	service.announceList(uuid.New(), hub, []*dto.DeviceResponse{{Name: "Batten"}})

	if events.published != 2 {
		t.Errorf("published %d times, want 2 — one for each distinct list", events.published)
	}
}

type recordingEvents struct {
	published int
}

func (r *recordingEvents) PublishJSON(_ uuid.UUID, _ any) {
	r.published++
}

func (r *recordingEvents) Publish(_ uuid.UUID, _ []byte) {}

func (r *recordingEvents) Subscribe(_ uuid.UUID) (<-chan []byte, func()) {
	return nil, func() {}
}

func (r *recordingEvents) Clients(_ uuid.UUID) int { return 0 }
