package device

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"time"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	devicerepo "ava/api/internal/repository/device"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"

	"github.com/google/uuid"
)

func (s *deviceService) SyncFromHub(ctx context.Context, tenantID, hubID uuid.UUID, req *dto.SyncDevicesRequest) ([]*dto.DeviceResponse, error) {
	now := time.Now()

	devices := make([]*model.Device, 0, len(req.Devices))

	for at := range req.Devices {
		item := &req.Devices[at]

		device := model.NewDevice(tenantID, hubID, &model.Reported{
			ExternalID:   item.ExternalID,
			Name:         item.Name,
			Kind:         item.Kind,
			Vendor:       item.Vendor,
			Model:        item.Model,
			IP:           item.IP,
			Parent:       item.Parent,
			Status:       model.DeviceStatus(item.Status),
			Capabilities: encode(item.Capabilities, "[]"),
			State:        encode(item.State, "{}"),
		})

		if device.Status == model.DeviceStatusOnline {
			device.LastSeenAt = &now
		}

		devices = append(devices, device)
	}

	if err := s.deviceRepo.SyncHubDevices(ctx, tenantID, hubID, devices); err != nil {
		logger.Error("device.SyncFromHub", logger.Err(err))

		return nil, err
	}

	synced, err := s.ListByHub(ctx, tenantID, hubID)
	if err != nil {
		return nil, err
	}

	s.announceList(tenantID, hubID, synced)

	return synced, nil
}

func (s *deviceService) MarkHubOffline(ctx context.Context, tenantID, hubID uuid.UUID) error {
	affected, err := s.deviceRepo.MarkHubDevicesOffline(ctx, hubID)
	if err != nil {
		logger.Error("device.MarkHubOffline", logger.Err(err))

		return err
	}

	if affected == 0 {
		return nil
	}

	devices, err := s.ListByHub(ctx, tenantID, hubID)
	if err != nil {
		return err
	}

	s.announceList(tenantID, hubID, devices)

	logger.Info("DEVICES_MARKED_OFFLINE",
		logger.Any("hub.ID", hubID),
		logger.Int("count", int(affected)),
	)

	return nil
}

// announceList tells clients what a hub holds only when it changed, so identical sweeps do not reset them.
func (s *deviceService) announceList(tenantID, hubID uuid.UUID, devices []*dto.DeviceResponse) {
	if s.events == nil {
		return
	}

	event := dto.NewDeviceListEvent(hubID, devices)

	encoded, err := json.Marshal(event)
	if err == nil && !s.listChanged(hubID, encoded) {
		return
	}

	s.events.PublishJSON(tenantID, event)
}

// listChanged records the newest list for a hub and says whether it differs.
func (s *deviceService) listChanged(hubID uuid.UUID, encoded []byte) bool {
	sum := sha256.Sum256(encoded)

	s.listsMu.Lock()
	defer s.listsMu.Unlock()

	if s.lists == nil {
		s.lists = make(map[uuid.UUID][32]byte)
	}

	if previous, seen := s.lists[hubID]; seen && previous == sum {
		return false
	}

	s.lists[hubID] = sum

	return true
}

func (s *deviceService) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.DeviceResponse, error) {
	devices, err := s.deviceRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		logger.Error("device.ListByTenant", logger.Err(err))

		return nil, err
	}

	return toDeviceResponses(devices), nil
}

func (s *deviceService) ListByHub(ctx context.Context, tenantID, hubID uuid.UUID) ([]*dto.DeviceResponse, error) {
	devices, err := s.deviceRepo.ListByHub(ctx, tenantID, hubID)
	if err != nil {
		logger.Error("device.ListByHub", logger.Err(err))

		return nil, err
	}

	return toDeviceResponses(devices), nil
}

func (s *deviceService) Update(
	ctx context.Context,
	tenantID, deviceID uuid.UUID,
	req *dto.UpdateDeviceRequest,
) (*dto.DeviceResponse, error) {
	fields := make(map[string]any, 3)

	if req.Name != nil {
		fields["name"] = *req.Name
	}

	if req.ClearRoom {
		fields["room_id"] = nil
	} else if req.RoomID != nil {
		fields["room_id"] = *req.RoomID
	}

	if req.Appliance != nil {
		fields["appliance"] = *req.Appliance
	}

	if len(fields) == 0 {
		return nil, ErrNothingToUpdate
	}

	if err := s.deviceRepo.Update(ctx, tenantID, deviceID, fields); err != nil {
		if serrors.Is(err, devicerepo.ErrDeviceNotFound) {
			return nil, ErrDeviceNotFound
		}

		logger.Error("device.Update", logger.Err(err))

		return nil, err
	}

	device, err := s.deviceRepo.GetByID(ctx, tenantID, deviceID)
	if err != nil {
		logger.Error("device.Update", logger.Err(err))

		return nil, err
	}

	return toDeviceResponse(device), nil
}

func encode(value any, fallback string) json.RawMessage {
	raw, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(fallback)
	}

	return raw
}
