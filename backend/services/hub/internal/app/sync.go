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
	entries := inventory.Discover(ctx, a.cfg.DiscoveryTimeout)

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

	items := inventory.ToSyncItems(entries)

	synced, err := a.client.SyncDevices(ctx, items)
	if err != nil {
		logger.Warn("DEVICE_SYNC_FAILED", logger.String("error", err.Error()))

		return
	}

	logger.Info("DEVICES_SYNCED",
		logger.Int("discovered", len(entries)),
		logger.Int("accepted", len(synced)),
	)
}

func (r *registry) get(externalID string) (device.Device, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	found, ok := r.devices[externalID]

	return found, ok
}
