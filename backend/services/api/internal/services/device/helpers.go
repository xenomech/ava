package device

import (
	"encoding/json"
	"fmt"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	"ava/pkg/logger"
	"ava/pkg/wire"
)

func toDeviceResponse(device *model.Device) *dto.DeviceResponse {
	capabilities, err := decodeCapabilities(device.Capabilities)
	if err != nil {
		logger.Warn("DEVICE_CAPABILITIES_UNREADABLE", logger.Any("device.ID", device.ID), logger.Err(err))
	}

	return &dto.DeviceResponse{
		ID:           device.ID,
		HubID:        device.HubID,
		ExternalID:   device.ExternalID,
		Name:         device.Name,
		RoomID:       device.RoomID,
		Room:         roomName(device.Room),
		Appliance:    device.Appliance,
		Kind:         device.Kind,
		Vendor:       device.Vendor,
		Model:        device.Model,
		Parent:       device.Parent,
		Status:       string(device.Status),
		LastSeenAt:   device.LastSeenAt,
		Capabilities: capabilities,
		State:        decodeState(device.State),
		CreatedAt:    device.CreatedAt,
	}
}

func decodeCapabilities(raw json.RawMessage) (wire.Capabilities, error) {
	if len(raw) == 0 {
		return wire.Capabilities{}, nil
	}

	var capabilities wire.Capabilities
	if err := json.Unmarshal(raw, &capabilities); err != nil {
		return wire.Capabilities{}, fmt.Errorf("decode capabilities: %w", err)
	}

	return capabilities, nil
}

func decodeState(raw json.RawMessage) wire.State {
	if len(raw) == 0 {
		return wire.State{}
	}

	var state wire.State
	if err := json.Unmarshal(raw, &state); err != nil {
		return wire.State{}
	}

	return state
}

func toDeviceResponses(devices []*model.Device) []*dto.DeviceResponse {
	responses := make([]*dto.DeviceResponse, 0, len(devices))
	for _, device := range devices {
		responses = append(responses, toDeviceResponse(device))
	}

	return responses
}

func roomName(room *model.Room) string {
	if room == nil {
		return ""
	}

	return room.Name
}
