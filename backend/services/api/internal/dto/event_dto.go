package dto

import "github.com/google/uuid"

const (
	EventDeviceState = "device.state"
	EventDeviceList  = "device.list"
	EventHubPresence = "hub.presence"
)

type DeviceStateEvent struct {
	Type   string          `json:"type"`
	Device *DeviceResponse `json:"device"`
}

func NewDeviceStateEvent(device *DeviceResponse) DeviceStateEvent {
	return DeviceStateEvent{Type: EventDeviceState, Device: device}
}

type DeviceListEvent struct {
	Type    string            `json:"type"`
	HubID   uuid.UUID         `json:"hub_id"`
	Devices []*DeviceResponse `json:"devices"`
}

func NewDeviceListEvent(hubID uuid.UUID, devices []*DeviceResponse) DeviceListEvent {
	if devices == nil {
		devices = []*DeviceResponse{}
	}

	return DeviceListEvent{Type: EventDeviceList, HubID: hubID, Devices: devices}
}

type HubPresenceEvent struct {
	Type   string    `json:"type"`
	HubID  uuid.UUID `json:"hub_id"`
	Online bool      `json:"online"`
}

func NewHubPresenceEvent(hubID uuid.UUID, online bool) HubPresenceEvent {
	return HubPresenceEvent{Type: EventHubPresence, HubID: hubID, Online: online}
}
