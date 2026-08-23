package app

import (
	"context"
	"encoding/json"
	"sync"

	"ava/pkg/logger"

	"ava/hub/internal/api"
	"ava/hub/internal/device"
	"ava/pkg/mqtt"
	"ava/pkg/wire"
)

func (a *App) startCommands(ctx context.Context, tokens *api.HubTokens) (*mqtt.Client, error) {
	topics := wire.TopicsFor(tokens.Tenant.Slug, tokens.Hub.ID)

	offline, err := json.Marshal(wire.Presence{Online: false, HubID: tokens.Hub.ID})
	if err != nil {
		return nil, err
	}

	online, err := json.Marshal(wire.Presence{Online: true, HubID: tokens.Hub.ID})
	if err != nil {
		return nil, err
	}

	client, err := mqtt.Connect(ctx, &mqtt.Options{
		BrokerURL: a.cfg.MQTTBrokerURL,
		ClientID:  "ava-hub-" + tokens.Hub.ID,
		Username:  a.state.BrokerUsername,
		Password:  a.state.BrokerPassword,
		WillTopic: topics.Status,
		Durable:   true,
		Will:      offline,
		OnConnect: func(client *mqtt.Client) {
			a.announce(ctx, client, topics, online)
		},
	})
	if err != nil {
		return nil, err
	}

	a.topics = topics
	a.mqtt = client

	return client, nil
}

func (a *App) announce(ctx context.Context, client *mqtt.Client, topics wire.Topics, online []byte) {
	if err := client.Subscribe(ctx, topics.Command, func(_ string, payload []byte) {
		a.handleCommand(ctx, payload)
	}); err != nil {
		logger.Warn("COMMAND_SUBSCRIBE_FAILED", logger.String("error", err.Error()))

		return
	}

	if err := client.Subscribe(ctx, topics.Apply, func(_ string, payload []byte) {
		a.handleApply(ctx, payload)
	}); err != nil {
		logger.Warn("APPLY_SUBSCRIBE_FAILED", logger.String("error", err.Error()))

		return
	}

	if err := client.Publish(ctx, topics.Status, online, true); err != nil {
		logger.Warn("PRESENCE_PUBLISH_FAILED", logger.String("error", err.Error()))

		return
	}

	logger.Info("HUB_ONLINE", logger.String("topic", topics.Status))
}

func (a *App) handleCommand(ctx context.Context, payload []byte) {
	cmd, err := wire.DecodeCommand(payload)
	if err != nil {
		logger.Warn("COMMAND_REJECTED", logger.String("error", err.Error()))

		return
	}

	target, ok := a.devices.get(cmd.DeviceID)
	if !ok {
		logger.Warn("COMMAND_UNKNOWN_DEVICE", logger.String("device_id", cmd.DeviceID))

		return
	}

	if err := target.Apply(ctx, cmd.Trait, cmd.Value); err != nil {
		logger.Warn("COMMAND_FAILED",
			logger.String("device_id", cmd.DeviceID),
			logger.String("trait", string(cmd.Trait)),
			logger.String("error", err.Error()),
		)

		return
	}

	logger.Info("COMMAND_APPLIED",
		logger.String("device_id", cmd.DeviceID),
		logger.String("trait", string(cmd.Trait)),
	)

	a.publishState(ctx, cmd.DeviceID, target)
}

func (a *App) publishState(ctx context.Context, deviceID string, target device.Device) {
	if a.mqtt == nil {
		return
	}

	state, err := target.State(ctx)
	if err != nil {
		logger.Warn("STATE_READ_FAILED", logger.String("device_id", deviceID), logger.String("error", err.Error()))

		return
	}

	payload, err := json.Marshal(wire.StateEvent{DeviceID: deviceID, State: state})
	if err != nil {
		return
	}

	if err := a.mqtt.Publish(ctx, a.topics.State, payload, false); err != nil {
		logger.Warn("STATE_PUBLISH_FAILED", logger.String("error", err.Error()))
	}
}

func (a *App) handleApply(ctx context.Context, payload []byte) {
	apply, err := wire.DecodeApply(payload)
	if err != nil {
		logger.Warn("APPLY_REJECTED", logger.String("error", err.Error()))

		return
	}

	var wg sync.WaitGroup

	for at := range apply.Targets {
		target := apply.Targets[at]

		handle, ok := a.devices.get(target.DeviceID)
		if !ok {
			logger.Warn("APPLY_UNKNOWN_DEVICE", logger.String("device_id", target.DeviceID))

			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := handle.Apply(ctx, target.Trait, target.Value); err != nil {
				logger.Warn("APPLY_FAILED",
					logger.String("device_id", target.DeviceID),
					logger.String("trait", string(target.Trait)),
					logger.String("error", err.Error()),
				)

				return
			}

			a.publishState(ctx, target.DeviceID, handle)
		}()
	}

	wg.Wait()

	logger.Info("APPLY_DONE", logger.Int("targets", len(apply.Targets)))
}
