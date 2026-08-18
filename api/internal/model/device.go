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
	TenantID   uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Tenant     *Tenant         `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	HubID      uuid.UUID       `gorm:"type:uuid;not null;index;uniqueIndex:idx_device_hub_external" json:"hub_id"`
	Hub        *Hub            `gorm:"foreignKey:HubID;constraint:OnDelete:CASCADE" json:"-"`
	ExternalID string          `gorm:"type:varchar(128);not null;uniqueIndex:idx_device_hub_external" json:"external_id"`
	Name       string          `gorm:"not null" json:"name"`
	Kind       string          `gorm:"type:varchar(40);not null" json:"kind"`
	Status     DeviceStatus    `gorm:"type:varchar(20);not null;default:'offline';index" json:"status"`
	LastSeenAt *time.Time      `json:"last_seen_at,omitempty"`
	State      json.RawMessage `gorm:"type:jsonb;default:'{}'" json:"state"`
}

func (device *Device) BeforeCreate(tx *gorm.DB) error {
	return device.BaseModel.BeforeCreate(tx)
}

func NewDevice(tenantID, hubID uuid.UUID, externalID, name, kind string, status DeviceStatus, state json.RawMessage) *Device {
	if len(state) == 0 {
		state = json.RawMessage("{}")
	}

	return &Device{
		TenantID:   tenantID,
		HubID:      hubID,
		ExternalID: externalID,
		Name:       name,
		Kind:       kind,
		Status:     status,
		State:      state,
	}
}
