package device

import (
	"context"
	"encoding/json"
	"fmt"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"
	"ava/pkg/wire"

	"github.com/google/uuid"
)

var ErrNothingToApply = serrors.NewCoded("nothing_to_apply", "no target could be applied")

func (s *deviceService) Apply(
	ctx context.Context,
	tenantID uuid.UUID,
	req *dto.ApplyRequest,
) (*dto.ApplyResponse, error) {
	if s.commander == nil {
		return nil, ErrCommandChannelUnavailable
	}

	wanted := make([]uuid.UUID, 0, len(req.Targets))
	for at := range req.Targets {
		wanted = append(wanted, req.Targets[at].DeviceID)
	}

	devices, err := s.deviceRepo.ListWithRelations(ctx, tenantID, wanted)
	if err != nil {
		logger.Error("device.Apply", logger.Err(err))

		return nil, err
	}

	known := make(map[uuid.UUID]*model.Device, len(devices))
	for _, found := range devices {
		known[found.ID] = found
	}

	out := &dto.ApplyResponse{Applied: []uuid.UUID{}, Skipped: []dto.SkippedTarget{}}
	batches := make(map[string][]wire.ApplyTarget)

	for at := range req.Targets {
		target := &req.Targets[at]

		accepted, reason := s.plan(known[target.DeviceID], target)
		if reason != "" {
			out.Skipped = append(out.Skipped, dto.SkippedTarget{DeviceID: target.DeviceID, Reason: reason})

			continue
		}

		topic := wire.TopicsFor(known[target.DeviceID].Tenant.Slug, known[target.DeviceID].HubID.String()).Apply
		batches[topic] = append(batches[topic], accepted)
		out.Applied = append(out.Applied, target.DeviceID)
	}

	if len(batches) == 0 {
		return out, ErrNothingToApply
	}

	if err := s.publishBatches(ctx, batches); err != nil {
		return nil, err
	}

	logger.Info("APPLY_PUBLISHED",
		logger.Int("hubs", len(batches)),
		logger.Int("applied", len(out.Applied)),
		logger.Int("skipped", len(out.Skipped)),
	)

	return out, nil
}

func (s *deviceService) plan(target *model.Device, wanted *dto.ApplyTargetRequest) (accepted wire.ApplyTarget, skip string) {
	if target == nil {
		return wire.ApplyTarget{}, "device not found"
	}

	if target.Hub == nil || target.Tenant == nil {
		return wire.ApplyTarget{}, "device not found"
	}

	if target.Hub.PresenceKnown() && !target.Hub.IsOnline() {
		return wire.ApplyTarget{}, "hub offline"
	}

	if target.Status == model.DeviceStatusOffline {
		return wire.ApplyTarget{}, "device offline"
	}

	capabilities, err := decodeCapabilities(target.Capabilities)
	if err != nil {
		return wire.ApplyTarget{}, "capabilities unreadable"
	}

	if err := capabilities.ValidateWrite(wanted.Trait, wanted.Value); err != nil {
		return wire.ApplyTarget{}, err.Error()
	}

	return wire.ApplyTarget{
		DeviceID: target.ExternalID,
		Trait:    wanted.Trait,
		Value:    wanted.Value,
	}, ""
}

func (s *deviceService) publishBatches(ctx context.Context, batches map[string][]wire.ApplyTarget) error {
	for topic, targets := range batches {
		payload, err := json.Marshal(wire.Apply{Targets: targets})
		if err != nil {
			return fmt.Errorf("encode apply: %w", err)
		}

		if err := s.commander.Publish(ctx, topic, payload, false); err != nil {
			logger.Error("APPLY_PUBLISH_FAILED", logger.String("topic", topic), logger.Err(err))

			return ErrCommandChannelUnavailable
		}
	}

	return nil
}
