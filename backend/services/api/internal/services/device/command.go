package device

import (
	"context"
	"encoding/json"
	"fmt"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	devicerepo "ava/api/internal/repository/device"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"
	"ava/pkg/wire"

	"github.com/google/uuid"
)

func (s *deviceService) SendCommand(
	ctx context.Context,
	tenantID, deviceID uuid.UUID,
	req *dto.SendCommandRequest,
) (*dto.CommandAcceptedResponse, error) {
	if s.commander == nil {
		return nil, ErrCommandChannelUnavailable
	}

	target, err := s.deviceRepo.GetByID(ctx, tenantID, deviceID)
	if err != nil {
		if serrors.Is(err, devicerepo.ErrDeviceNotFound) {
			return nil, ErrDeviceNotFound
		}

		return nil, err
	}

	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(wire.Command{
		DeviceID: target.ExternalID,
		Action:   req.Action,
		Value:    req.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("encode command: %w", err)
	}

	topic := wire.TopicsFor(tenant.Slug, target.HubID.String()).Command

	if err := s.commander.Publish(ctx, topic, payload, false); err != nil {
		logger.Error("COMMAND_PUBLISH_FAILED", logger.String("topic", topic), logger.Err(err))

		return nil, ErrCommandChannelUnavailable
	}

	logger.Info("COMMAND_PUBLISHED",
		logger.String("topic", topic),
		logger.String("device.ExternalID", target.ExternalID),
		logger.String("action", req.Action),
	)

	return &dto.CommandAcceptedResponse{
		DeviceID:   target.ID,
		ExternalID: target.ExternalID,
		Action:     req.Action,
		Topic:      topic,
	}, nil
}

var ErrCommandChannelUnavailable = serrors.NewCoded("command_channel_unavailable", "the hub command channel is unavailable")

func (s *deviceService) ApplyReportedState(
	ctx context.Context,
	hubID uuid.UUID,
	externalID string,
	state json.RawMessage,
) error {
	updated, err := s.deviceRepo.ApplyState(ctx, hubID, externalID, state)
	if err != nil {
		if serrors.Is(err, devicerepo.ErrDeviceNotFound) {
			return ErrDeviceNotFound
		}

		return err
	}

	s.announce(updated)

	return nil
}

func (s *deviceService) announce(device *model.Device) {
	if s.events == nil {
		return
	}

	s.events.PublishJSON(device.TenantID, map[string]any{
		"type":   "device.state",
		"device": toDeviceResponse(device),
	})
}
