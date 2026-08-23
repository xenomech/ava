package inventory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ava/pkg/logger"
	"ava/pkg/wire"

	"ava/hub/external/wiz"
	"ava/hub/internal/device"
)

const (
	KindBulb = "bulb"
	KindPlug = "plug"

	StatusOnline = "online"
)

type Entry struct {
	Spec  device.Spec
	State wire.State
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
		entries = append(entries, describeWiz(ctx, &lights[at], timeout))
	}

	return entries
}

func describeWiz(ctx context.Context, found *wiz.Found, timeout time.Duration) Entry {
	entry := Entry{
		Spec: device.Spec{
			Vendor: device.VendorWiz,
			ID:     found.Info.ID,
			IP:     found.Info.IP,
		},
		State: found.State,
		Kind:  KindBulb,
		Name:  friendlyName("WiZ", found.Info.ID),
	}

	light := wiz.New(found.Info.IP, timeout)

	info, err := light.Identify(ctx)
	if err != nil {
		logger.Warn("WIZ_IDENTIFY_FAILED",
			logger.String("ip", found.Info.IP),
			logger.Err(err),
		)

		return entry
	}

	entry.Spec.Capabilities = light.Capabilities()
	entry.Kind = kindFor(entry.Spec.Capabilities)
	entry.Name = friendlyName("WiZ", info.MAC)

	if info.Model != "" {
		entry.Spec.Name = info.Model
	}

	return entry
}

func kindFor(capabilities wire.Capabilities) string {
	if capabilities.Has(wire.TraitBrightness) {
		return KindBulb
	}

	return KindPlug
}

func friendlyName(vendor, id string) string {
	trimmed := strings.ReplaceAll(id, ":", "")
	if len(trimmed) > 4 {
		trimmed = trimmed[len(trimmed)-4:]
	}

	return fmt.Sprintf("%s %s", vendor, strings.ToUpper(trimmed))
}

func ToSyncItems(entries []Entry) []wire.DeviceReport {
	items := make([]wire.DeviceReport, 0, len(entries))

	for at := range entries {
		entry := &entries[at]

		items = append(items, wire.DeviceReport{
			ExternalID:   entry.Spec.ID,
			Name:         entry.Name,
			Kind:         entry.Kind,
			Status:       StatusOnline,
			Vendor:       string(entry.Spec.Vendor),
			Model:        entry.Spec.Name,
			IP:           entry.Spec.IP,
			Capabilities: entry.Spec.Capabilities,
			State:        entry.State,
		})
	}

	return items
}
