package app

import (
	"context"
	"sync"
	"time"

	"ava/pkg/logger"

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

	logger.Info("DEVICES_SYNCED",
		logger.Int("reported", len(entries)),
		logger.Int("known", len(synced)),
	)
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
