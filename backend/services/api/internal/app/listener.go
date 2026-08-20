package app

import (
	"context"
	"encoding/json"
	"strings"

	"ava/pkg/logger"
	"ava/pkg/wire"

	"github.com/google/uuid"
)

const (
	stateTopicFilter    = "ava/+/+/state"
	presenceTopicFilter = "ava/+/+/status"
)

func (a *App) listenForState(ctx context.Context) {
	if a.Publisher == nil {
		return
	}

	err := a.Publisher.Subscribe(ctx, stateTopicFilter, func(topic string, payload []byte) {
		a.applyState(ctx, topic, payload)
	})
	if err != nil {
		logger.Warn("STATE_SUBSCRIBE_FAILED", logger.Err(err))
	}
}

func (a *App) applyState(ctx context.Context, topic string, payload []byte) {
	hubID, ok := hubFromTopic(topic)
	if !ok {
		return
	}

	var event wire.StateEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		logger.Warn("STATE_DECODE_FAILED", logger.String("topic", topic))

		return
	}

	if event.DeviceID == "" {
		return
	}

	patch, err := json.Marshal(map[string]any{
		"power":      event.Power,
		"brightness": event.Brightness,
		"color_temp": event.ColorTemp,
	})
	if err != nil {
		return
	}

	if err := a.Service.Device.ApplyReportedState(ctx, hubID, event.DeviceID, patch); err != nil {
		logger.Warn("STATE_APPLY_FAILED",
			logger.String("device.ExternalID", event.DeviceID),
			logger.Err(err),
		)

		return
	}

	logger.Debug("STATE_APPLIED", logger.String("device.ExternalID", event.DeviceID))
}

func (a *App) listenForPresence(ctx context.Context) {
	if a.Publisher == nil {
		return
	}

	err := a.Publisher.Subscribe(ctx, presenceTopicFilter, func(topic string, payload []byte) {
		a.applyPresence(ctx, topic, payload)
	})
	if err != nil {
		logger.Warn("PRESENCE_SUBSCRIBE_FAILED", logger.Err(err))
	}
}

func (a *App) applyPresence(ctx context.Context, topic string, payload []byte) {
	hubID, ok := hubFromTopic(topic)
	if !ok {
		return
	}

	var presence wire.Presence

	if err := json.Unmarshal(payload, &presence); err != nil {
		logger.Warn("PRESENCE_DECODE_FAILED", logger.String("topic", topic))

		return
	}

	hub, err := a.Service.Hub.ApplyPresence(ctx, hubID, presence.Online)
	if err != nil {
		logger.Warn("PRESENCE_APPLY_FAILED", logger.Any("hub.ID", hubID), logger.Err(err))

		return
	}

	if hub.IsOnline() {
		return
	}

	if err := a.Service.Device.MarkHubOffline(ctx, hub.TenantID, hub.ID); err != nil {
		logger.Warn("PRESENCE_CASCADE_FAILED", logger.Any("hub.ID", hubID), logger.Err(err))
	}
}

func hubFromTopic(topic string) (uuid.UUID, bool) {
	parts := strings.Split(topic, "/")
	if len(parts) != 4 {
		return uuid.Nil, false
	}

	hubID, err := uuid.Parse(parts[2])
	if err != nil {
		return uuid.Nil, false
	}

	return hubID, true
}
