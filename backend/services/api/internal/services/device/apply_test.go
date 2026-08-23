package device

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	"ava/pkg/wire"

	"github.com/google/uuid"
)

const lightCapabilities = `[
	{"trait":"power","kind":"bool","access":"rw"},
	{"trait":"brightness","kind":"number","access":"rw","min":10,"max":100,"unit":"%"},
	{"trait":"temperature","kind":"number","access":"r","unit":"C"}
]`

func onlineHub() *model.Hub {
	seen := time.Now()

	return &model.Hub{Online: true, LastSeenAt: &seen}
}

func lamp(hub *model.Hub, status model.DeviceStatus, capabilities string) *model.Device {
	return &model.Device{
		ExternalID:   "lamp-1",
		Status:       status,
		Hub:          hub,
		Tenant:       &model.Tenant{Slug: "acme"},
		Capabilities: json.RawMessage(capabilities),
	}
}

func want(trait wire.Trait, value wire.Value) *dto.ApplyTargetRequest {
	return &dto.ApplyTargetRequest{DeviceID: uuid.New(), Trait: trait, Value: value}
}

func TestATargetIsAcceptedWhenEverythingLinesUp(t *testing.T) {
	service := &deviceService{}

	accepted, skip := service.plan(
		lamp(onlineHub(), model.DeviceStatusOnline, lightCapabilities),
		want(wire.TraitBrightness, wire.Number(60)),
	)

	if skip != "" {
		t.Fatalf("skipped: %s", skip)
	}

	if accepted.DeviceID != "lamp-1" {
		t.Errorf("device id = %s, want the external id", accepted.DeviceID)
	}

	if number, _ := accepted.Value.Number(); number != 60 {
		t.Errorf("value = %v", accepted.Value)
	}
}

func TestATargetIsSkippedWithAReasonTheCallerCanShow(t *testing.T) {
	stale := time.Now().Add(-time.Hour)
	offlineHub := &model.Hub{Online: false, LastSeenAt: &stale}

	cases := map[string]struct {
		device *model.Device
		wanted *dto.ApplyTargetRequest
		reason string
	}{
		"device not in the tenant": {
			nil, want(wire.TraitPower, wire.Bool(true)), "device not found",
		},
		"relations not loaded": {
			&model.Device{ExternalID: "lamp-1"}, want(wire.TraitPower, wire.Bool(true)), "device not found",
		},
		"hub offline": {
			lamp(offlineHub, model.DeviceStatusOnline, lightCapabilities),
			want(wire.TraitPower, wire.Bool(true)), "hub offline",
		},
		"device offline": {
			lamp(onlineHub(), model.DeviceStatusOffline, lightCapabilities),
			want(wire.TraitPower, wire.Bool(true)), "device offline",
		},
		"capabilities unreadable": {
			lamp(onlineHub(), model.DeviceStatusOnline, `{"not":"an array"}`),
			want(wire.TraitPower, wire.Bool(true)), "capabilities unreadable",
		},
		"trait the device does not have": {
			lamp(onlineHub(), model.DeviceStatusOnline, lightCapabilities),
			want(wire.TraitColorTemp, wire.Number(3000)), "does not have this trait",
		},
		"read only trait": {
			lamp(onlineHub(), model.DeviceStatusOnline, lightCapabilities),
			want(wire.TraitTemperature, wire.Number(20)), "read only",
		},
		"value out of range": {
			lamp(onlineHub(), model.DeviceStatusOnline, lightCapabilities),
			want(wire.TraitBrightness, wire.Number(140)), "out of range",
		},
		"value of the wrong type": {
			lamp(onlineHub(), model.DeviceStatusOnline, lightCapabilities),
			want(wire.TraitBrightness, wire.Bool(true)), "wrong type",
		},
	}

	service := &deviceService{}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			accepted, skip := service.plan(tc.device, tc.wanted)

			if !strings.Contains(skip, tc.reason) {
				t.Errorf("reason = %q, want it to mention %q", skip, tc.reason)
			}

			if accepted.DeviceID != "" {
				t.Errorf("a skipped target must not be published, got %+v", accepted)
			}
		})
	}
}

func TestAHubThatHasNeverReportedIsGivenTheBenefitOfTheDoubt(t *testing.T) {
	service := &deviceService{}

	_, skip := service.plan(
		lamp(&model.Hub{}, model.DeviceStatusOnline, lightCapabilities),
		want(wire.TraitPower, wire.Bool(true)),
	)

	if skip != "" {
		t.Errorf("a freshly paired hub was skipped: %s", skip)
	}
}
