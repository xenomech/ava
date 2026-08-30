package dto

import (
	"time"

	"ava/pkg/wire"

	"github.com/google/uuid"
)

type SyncDeviceItem struct {
	wire.DeviceReport
}

type SyncDevicesRequest struct {
	Devices []SyncDeviceItem `json:"devices" validate:"required,dive"`
}

type UpdateDeviceRequest struct {
	Name      *string    `json:"name,omitempty" validate:"omitempty,max=100"`
	RoomID    *uuid.UUID `json:"room_id,omitempty"`
	ClearRoom bool       `json:"clear_room,omitempty"`
	Appliance *string    `json:"appliance,omitempty" validate:"omitempty,max=40"`
}

type DeviceResponse struct {
	ID           uuid.UUID         `json:"id"`
	HubID        uuid.UUID         `json:"hub_id"`
	ExternalID   string            `json:"external_id"`
	Name         string            `json:"name"`
	RoomID       *uuid.UUID        `json:"room_id,omitempty"`
	Room         string            `json:"room"`
	Appliance    string            `json:"appliance"`
	Kind         string            `json:"kind"`
	Vendor       string            `json:"vendor,omitempty"`
	Model        string            `json:"model,omitempty"`
	Parent       string            `json:"parent,omitempty"`
	Status       string            `json:"status"`
	LastSeenAt   *time.Time        `json:"last_seen_at,omitempty"`
	Capabilities wire.Capabilities `json:"capabilities"`
	State        wire.State        `json:"state"`
	CreatedAt    time.Time         `json:"created_at"`
}

type SendCommandRequest struct {
	Trait wire.Trait `json:"trait" validate:"required,max=64"`
	Value wire.Value `json:"value"`
}

type CommandAcceptedResponse struct {
	DeviceID   uuid.UUID  `json:"device_id"`
	ExternalID string     `json:"external_id"`
	Trait      wire.Trait `json:"trait"`
	Topic      string     `json:"topic"`
}

type ApplyTargetRequest struct {
	DeviceID uuid.UUID  `json:"device_id" validate:"required"`
	Trait    wire.Trait `json:"trait" validate:"required,max=64"`
	Value    wire.Value `json:"value"`
}

type ApplyRequest struct {
	Targets []ApplyTargetRequest `json:"targets" validate:"required,min=1,max=100,dive"`
}

type SkippedTarget struct {
	DeviceID uuid.UUID `json:"device_id"`
	Reason   string    `json:"reason"`
}

type ApplyResponse struct {
	Applied []uuid.UUID     `json:"applied"`
	Skipped []SkippedTarget `json:"skipped"`
}
