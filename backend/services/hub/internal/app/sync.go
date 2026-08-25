package app

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"ava/pkg/logger"
	"ava/pkg/wire"

	"ava/hub/internal/device"
	"ava/hub/internal/device/adapters"
	"ava/hub/internal/inventory"
)

type registry struct {
	mu      sync.RWMutex
	devices map[string]device.Device
}

func newRegistry() *registry {
	return &registry{devices: make(map[string]device.Device)}
}

func (r *registry) replace(devices map[string]device.Device) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.devices = devices
}

func (a *App) syncLoop(ctx context.Context) {
	ticker := time.NewTicker(a.cfg.SyncInterval)
	defer ticker.Stop()

	for {
		a.syncOnce(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (a *App) syncOnce(ctx context.Context) {
	entries := a.recoverMissing(ctx, inventory.Discover(ctx, a.cfg.DiscoveryTimeout))

	opened := make(map[string]device.Device, len(entries))

	for at := range entries {
		spec := entries[at].Spec

		handle, err := adapters.Open(&spec)
		if err != nil {
			logger.Warn("DEVICE_OPEN_FAILED",
				logger.String("device_id", spec.ID),
				logger.String("error", err.Error()),
			)

			continue
		}

		opened[spec.ID] = handle
	}

	a.devices.replace(opened)

	a.lastSeen = entries

	items := inventory.ToSyncItems(entries)

	synced, err := a.client.SyncDevices(ctx, items)
	if err != nil {
		logger.Warn("DEVICE_SYNC_FAILED", logger.String("error", err.Error()))

		return
	}

	a.publishSweep(ctx, entries)

	logger.Info("DEVICES_SYNCED",
		logger.Int("reported", len(entries)),
		logger.Int("known", len(synced)),
	)
}

// publishSweep reports what the sweep found, over the channel that merges.
//
// State used to travel with the sync itself, which replaced a device's stored
// reading outright — and a sweep spends several seconds discovering before it
// reports, so a snapshot taken six seconds ago landed on top of a change made
// one second ago. The light was right and the app said otherwise, which is the
// worst way round for it to be wrong.
//
// Sending it here instead gives state exactly one writer, and that writer
// merges: a reading can correct a trait or retire it, but it cannot undo
// something newer than itself.
func (a *App) publishSweep(ctx context.Context, entries []inventory.Entry) {
	if a.mqtt == nil || a.topics.State == "" {
		return
	}

	for at := range entries {
		entry := &entries[at]

		payload, err := json.Marshal(wire.StateEvent{
			DeviceID: entry.Spec.ID,
			State:    entry.State,
		})
		if err != nil {
			continue
		}

		if err := a.mqtt.Publish(ctx, a.topics.State, payload, false); err != nil {
			logger.Warn("SWEEP_PUBLISH_FAILED", logger.String("error", err.Error()))

			return
		}
	}
}

func (r *registry) get(externalID string) (device.Device, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	found, ok := r.devices[externalID]

	return found, ok
}

func (a *App) recoverMissing(ctx context.Context, found []inventory.Entry) []inventory.Entry {
	answered := make(map[string]struct{}, len(found))
	for at := range found {
		answered[found[at].Spec.ID] = struct{}{}
	}

	for at := range a.lastSeen {
		previous := a.lastSeen[at]

		if _, ok := answered[previous.Spec.ID]; ok {
			continue
		}

		spec := previous.Spec

		handle, err := adapters.Open(&spec)
		if err != nil {
			continue
		}

		state, err := handle.State(ctx)
		if err != nil {
			logger.Info("DEVICE_UNREACHABLE",
				logger.String("device_id", previous.Spec.ID),
				logger.String("ip", previous.Spec.IP),
			)

			continue
		}

		logger.Debug("DEVICE_RECOVERED", logger.String("device_id", previous.Spec.ID))

		previous.State = state
		found = append(found, previous)
	}

	return found
}
