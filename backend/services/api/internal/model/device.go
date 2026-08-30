package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusOffline DeviceStatus = "offline"
)

func (status DeviceStatus) IsValid() bool {
	switch status {
	case DeviceStatusOnline, DeviceStatusOffline:
		return true
	default:
		return false
	}
}

type Device struct {
	BaseModel
	TenantID     uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Tenant       *Tenant         `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	HubID        uuid.UUID       `gorm:"type:uuid;not null;index;uniqueIndex:idx_device_hub_external" json:"hub_id"`
	Hub          *Hub            `gorm:"foreignKey:HubID;constraint:OnDelete:CASCADE" json:"-"`
	ExternalID   string          `gorm:"type:varchar(128);not null;uniqueIndex:idx_device_hub_external" json:"external_id"`
	Name         string          `gorm:"not null" json:"name"`
	RoomID       *uuid.UUID      `gorm:"type:uuid;index" json:"room_id,omitempty"`
	Room         *Room           `gorm:"foreignKey:RoomID;constraint:OnDelete:SET NULL" json:"room,omitempty"`
	Appliance    string          `gorm:"type:varchar(40);not null;default:''" json:"appliance"`
	Kind         string          `gorm:"type:varchar(40);not null" json:"kind"`
	Vendor       string          `gorm:"type:varchar(40);not null;default:''" json:"vendor"`
	Model        string          `gorm:"type:varchar(80);not null;default:''" json:"model"`
	IP           string          `gorm:"type:varchar(45);not null;default:''" json:"ip"`
	Parent       string          `gorm:"type:varchar(128);not null;default:''" json:"parent"`
	Status       DeviceStatus    `gorm:"type:varchar(20);not null;default:'offline';index" json:"status"`
	LastSeenAt   *time.Time      `json:"last_seen_at,omitempty"`
	Capabilities json.RawMessage `gorm:"type:jsonb;default:'[]'" json:"capabilities"`
	State        json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"state"`
}

func (device *Device) BeforeCreate(tx *gorm.DB) error {
	return device.BaseModel.BeforeCreate(tx)
}

type Reported struct {
	ExternalID   string
	Name         string
	Kind         string
	Vendor       string
	Model        string
	IP           string
	Parent       string
	Status       DeviceStatus
	Capabilities json.RawMessage
	State        json.RawMessage
}

func NewDevice(tenantID, hubID uuid.UUID, reported *Reported) *Device {
	state := reported.State
	if len(state) == 0 {
		state = json.RawMessage("{}")
	}

	capabilities := reported.Capabilities
	if len(capabilities) == 0 {
		capabilities = json.RawMessage("[]")
	}

	return &Device{
		TenantID:     tenantID,
		HubID:        hubID,
		ExternalID:   reported.ExternalID,
		Name:         reported.Name,
		Kind:         reported.Kind,
		Vendor:       reported.Vendor,
		Model:        reported.Model,
		IP:           reported.IP,
		Parent:       reported.Parent,
		Status:       reported.Status,
		Capabilities: capabilities,
		State:        state,
	}
}
