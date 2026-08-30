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

	target, err := s.deviceRepo.GetWithRelations(ctx, tenantID, deviceID)
	if err != nil {
		if serrors.Is(err, devicerepo.ErrDeviceNotFound) {
			return nil, ErrDeviceNotFound
		}

		return nil, err
	}

	if target.Hub == nil || target.Tenant == nil {
		return nil, ErrDeviceNotFound
	}

	if target.Hub.PresenceKnown() && !target.Hub.IsOnline() {
		return nil, ErrHubOffline
	}

	capabilities, err := decodeCapabilities(target.Capabilities)
	if err != nil {
		return nil, err
	}

	if err := capabilities.ValidateWrite(req.Trait, req.Value); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrTraitRejected, err)
	}

	payload, err := json.Marshal(wire.Command{
		DeviceID: target.ExternalID,
		Trait:    req.Trait,
		Value:    req.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("encode command: %w", err)
	}

	topic := wire.TopicsFor(target.Tenant.Slug, target.HubID.String()).Command

	if err := s.commander.Publish(ctx, topic, payload, false); err != nil {
		logger.Error("COMMAND_PUBLISH_FAILED", logger.String("topic", topic), logger.Err(err))

		return nil, ErrCommandChannelUnavailable
	}

	logger.Info("COMMAND_PUBLISHED",
		logger.String("topic", topic),
		logger.String("device.ExternalID", target.ExternalID),
		logger.String("trait", string(req.Trait)),
	)

	return &dto.CommandAcceptedResponse{
		DeviceID:   target.ID,
		ExternalID: target.ExternalID,
		Trait:      req.Trait,
		Topic:      topic,
	}, nil
}

var (
	ErrCommandChannelUnavailable = serrors.NewCoded("command_channel_unavailable", "the hub command channel is unavailable")
	ErrHubOffline                = serrors.NewCoded("hub_offline", "the hub for this device is offline")
	ErrTraitRejected             = serrors.NewCoded("trait_rejected", "the device cannot accept this value")
)

func (s *deviceService) ApplyReportedState(
	ctx context.Context,
	hubID uuid.UUID,
	externalID string,
	state wire.State,
) error {
	// A trait reported as null is a retraction, not a value, so it must never be stored verbatim.
	set, cleared := state.Settled()

	patch, err := json.Marshal(set)
	if err != nil {
		return fmt.Errorf("encode reported state: %w", err)
	}

	retire := make([]string, 0, len(cleared))
	for _, trait := range cleared {
		retire = append(retire, string(trait))
	}

	updated, err := s.deviceRepo.ApplyState(ctx, hubID, externalID, patch, retire)
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

	s.events.PublishJSON(device.TenantID, dto.NewDeviceStateEvent(toDeviceResponse(device)))
}
