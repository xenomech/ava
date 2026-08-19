package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SyncDeviceItem struct {
	ExternalID string          `json:"external_id" validate:"required,max=128"`
	Name       string          `json:"name" validate:"required,max=100"`
	Kind       string          `json:"kind" validate:"required,max=40"`
	Status     string          `json:"status" validate:"required,oneof=online offline"`
	State      json.RawMessage `json:"state,omitempty"`
}

type SyncDevicesRequest struct {
	Devices []SyncDeviceItem `json:"devices" validate:"required,dive"`
}

type RenameDeviceRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type DeviceResponse struct {
	ID         uuid.UUID       `json:"id"`
	HubID      uuid.UUID       `json:"hub_id"`
	ExternalID string          `json:"external_id"`
	Name       string          `json:"name"`
	Kind       string          `json:"kind"`
	Status     string          `json:"status"`
	LastSeenAt *time.Time      `json:"last_seen_at,omitempty"`
	State      json.RawMessage `json:"state"`
	CreatedAt  time.Time       `json:"created_at"`
}

type SendCommandRequest struct {
	Action string          `json:"action" validate:"required,oneof=power brightness color_temp"`
	Value  json.RawMessage `json:"value" validate:"required"`
}

type CommandAcceptedResponse struct {
	DeviceID   uuid.UUID `json:"device_id"`
	ExternalID string    `json:"external_id"`
	Action     string    `json:"action"`
	Topic      string    `json:"topic"`
}
