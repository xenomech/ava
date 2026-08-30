package broker

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

	topicSegments = 4
	hubSegment    = 2
)

type DeviceStateHandler func(ctx context.Context, hubID uuid.UUID, externalID string, state wire.State) error

type HubPresenceHandler func(ctx context.Context, hubID uuid.UUID, online bool) error

func (b *Broker) Listen(ctx context.Context, onState DeviceStateHandler, onPresence HubPresenceHandler) {
	if b == nil {
		return
	}

	b.watch(ctx, stateTopicFilter, "STATE_SUBSCRIBE_FAILED", func(topic string, payload []byte) {
		b.deviceState(ctx, onState, topic, payload)
	})

	b.watch(ctx, presenceTopicFilter, "PRESENCE_SUBSCRIBE_FAILED", func(topic string, payload []byte) {
		b.hubPresence(ctx, onPresence, topic, payload)
	})
}

func (b *Broker) watch(ctx context.Context, filter, failure string, handle func(string, []byte)) {
	if err := b.client.Subscribe(ctx, filter, handle); err != nil {
		logger.Warn(failure, logger.Err(err))
	}
}

func (b *Broker) deviceState(ctx context.Context, apply DeviceStateHandler, topic string, payload []byte) {
	hubID, ok := hubFromTopic(topic)
	if !ok {
		return
	}

	event, err := wire.DecodeStateEvent(payload)
	if err != nil {
		logger.Warn("STATE_DECODE_FAILED", logger.String("topic", topic), logger.Err(err))

		return
	}

	if err := apply(ctx, hubID, event.DeviceID, event.State); err != nil {
		logger.Warn("STATE_APPLY_FAILED",
			logger.String("device.ExternalID", event.DeviceID),
			logger.Err(err),
		)

		return
	}

	logger.Debug("STATE_APPLIED", logger.String("device.ExternalID", event.DeviceID))
}

func (b *Broker) hubPresence(ctx context.Context, apply HubPresenceHandler, topic string, payload []byte) {
	hubID, ok := hubFromTopic(topic)
	if !ok {
		return
	}

	var presence wire.Presence

	if err := json.Unmarshal(payload, &presence); err != nil {
		logger.Warn("PRESENCE_DECODE_FAILED", logger.String("topic", topic))

		return
	}

	if err := apply(ctx, hubID, presence.Online); err != nil {
		logger.Warn("PRESENCE_APPLY_FAILED", logger.Any("hub.ID", hubID), logger.Err(err))
	}
}

func hubFromTopic(topic string) (uuid.UUID, bool) {
	parts := strings.Split(topic, "/")
	if len(parts) != topicSegments {
		return uuid.Nil, false
	}

	hubID, err := uuid.Parse(parts[hubSegment])
	if err != nil {
		return uuid.Nil, false
	}

	return hubID, true
}
