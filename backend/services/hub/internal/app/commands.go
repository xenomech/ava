package app

import (
	"context"
	"encoding/json"

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

	if err := apply(ctx, target, cmd); err != nil {
		logger.Warn("COMMAND_FAILED",
			logger.String("device_id", cmd.DeviceID),
			logger.String("action", cmd.Action),
			logger.String("error", err.Error()),
		)

		return
	}

	logger.Info("COMMAND_APPLIED",
		logger.String("device_id", cmd.DeviceID),
		logger.String("action", cmd.Action),
	)

	a.publishState(ctx, cmd.DeviceID, target)
}

func apply(ctx context.Context, target device.Device, cmd wire.Command) error {
	switch cmd.Action {
	case wire.ActionPower:
		on, err := cmd.Bool()
		if err != nil {
			return err
		}

		return target.SetPower(ctx, on)
	case wire.ActionBrightness:
		level, err := cmd.Int()
		if err != nil {
			return err
		}

		return target.SetBrightness(ctx, level)
	case wire.ActionColorTemp:
		kelvin, err := cmd.Int()
		if err != nil {
			return err
		}

		return target.SetColorTemp(ctx, kelvin)
	default:
		return wire.ErrUnknownAction
	}
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

	payload, err := json.Marshal(wire.StateEvent{
		DeviceID:   deviceID,
		Power:      state.Power,
		Brightness: state.Brightness,
		ColorTemp:  state.ColorTemp,
	})
	if err != nil {
		return
	}

	if err := a.mqtt.Publish(ctx, a.topics.State, payload, false); err != nil {
		logger.Warn("STATE_PUBLISH_FAILED", logger.String("error", err.Error()))
	}
}
