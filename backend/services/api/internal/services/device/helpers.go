package device

import (
	"encoding/json"

	"ava/api/internal/dto"
	"ava/api/internal/model"
)

func toDeviceResponse(device *model.Device) *dto.DeviceResponse {
	state := device.State
	if len(state) == 0 {
		state = json.RawMessage("{}")
	}

	return &dto.DeviceResponse{
		ID:         device.ID,
		HubID:      device.HubID,
		ExternalID: device.ExternalID,
		Name:       device.Name,
		Kind:       device.Kind,
		Status:     string(device.Status),
		LastSeenAt: device.LastSeenAt,
		State:      state,
		CreatedAt:  device.CreatedAt,
	}
}

func toDeviceResponses(devices []*model.Device) []*dto.DeviceResponse {
	responses := make([]*dto.DeviceResponse, 0, len(devices))
	for _, device := range devices {
		responses = append(responses, toDeviceResponse(device))
	}

	return responses
}
