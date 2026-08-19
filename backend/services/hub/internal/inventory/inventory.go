package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ava/pkg/logger"

	"ava/hub/external/tuya"
	"ava/hub/external/wiz"
	"ava/hub/internal/api"
	"ava/hub/internal/device"
)

const (
	KindBulb = "bulb"
	KindPlug = "plug"

	StatusOnline = "online"
)

type Entry struct {
	Spec  device.Spec
	State device.State
	Kind  string
	Name  string
}

func Discover(ctx context.Context, timeout time.Duration) []Entry {
	entries := make([]Entry, 0)

	lights, err := wiz.Discover(ctx, timeout)
	if err != nil {
		logger.Warn("WIZ_DISCOVERY_FAILED", logger.String("error", err.Error()))
	}

	for at := range lights {
		entries = append(entries, fromWiz(&lights[at]))
	}

	plugs, err := tuya.Discover(ctx, timeout)
	if err != nil {
		logger.Warn("TUYA_DISCOVERY_FAILED", logger.String("error", err.Error()))
	}

	for at := range plugs {
		if !plugs[at].Supported() {
			logger.Warn("TUYA_PROTOCOL_UNSUPPORTED",
				logger.String("device_id", plugs[at].Info.ID),
				logger.String("ip", plugs[at].Info.IP),
				logger.String("version", plugs[at].Version),
			)

			continue
		}

		entries = append(entries, fromTuya(&plugs[at]))
	}

	return entries
}

func fromWiz(found *wiz.Found) Entry {
	return Entry{
		Spec: device.Spec{
			Vendor:       device.VendorWiz,
			ID:           found.Info.ID,
			IP:           found.Info.IP,
			Capabilities: device.CapabilityBrightness | device.CapabilityColorTemp,
		},
		State: found.State,
		Kind:  KindBulb,
		Name:  friendlyName("WiZ", found.Info.ID),
	}
}

func fromTuya(found *tuya.Found) Entry {
	return Entry{
		Spec: device.Spec{
			Vendor: device.VendorTuya,
			ID:     found.Info.ID,
			IP:     found.Info.IP,
		},
		Kind: KindPlug,
		Name: friendlyName("Wipro", found.Info.ID),
	}
}

func friendlyName(vendor, id string) string {
	trimmed := strings.ReplaceAll(id, ":", "")
	if len(trimmed) > 4 {
		trimmed = trimmed[len(trimmed)-4:]
	}

	return fmt.Sprintf("%s %s", vendor, strings.ToUpper(trimmed))
}

func ToSyncItems(entries []Entry) []api.SyncDeviceItem {
	items := make([]api.SyncDeviceItem, 0, len(entries))

	for at := range entries {
		entry := &entries[at]

		payload := statePayload{
			Power:        entry.State.Power,
			Brightness:   entry.State.Brightness,
			ColorTemp:    entry.State.ColorTemp,
			Capabilities: entry.Spec.Capabilities.Names(),
			Vendor:       string(entry.Spec.Vendor),
			IP:           entry.Spec.IP,
		}

		raw, err := json.Marshal(payload)
		if err != nil {
			logger.Warn("DEVICE_STATE_ENCODE_FAILED", logger.String("device_id", entry.Spec.ID))

			continue
		}

		items = append(items, api.SyncDeviceItem{
			ExternalID: entry.Spec.ID,
			Name:       entry.Name,
			Kind:       entry.Kind,
			Status:     StatusOnline,
			State:      raw,
		})
	}

	return items
}

type statePayload struct {
	Power        bool     `json:"power"`
	Brightness   int      `json:"brightness,omitempty"`
	ColorTemp    int      `json:"color_temp,omitempty"`
	Capabilities []string `json:"capabilities"`
	Vendor       string   `json:"vendor"`
	IP           string   `json:"ip,omitempty"`
}
