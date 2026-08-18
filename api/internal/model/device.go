package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceStatus string

const (
	DeviceStatusActive  DeviceStatus = "active"
	DeviceStatusRevoked DeviceStatus = "revoked"
)

type Device struct {
	BaseModel
	TenantID     uuid.UUID    `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Tenant       *Tenant      `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	Name         string       `gorm:"not null" json:"name"`
	Status       DeviceStatus `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	RefreshToken string       `gorm:"type:varchar(64);not null;uniqueIndex" json:"-"`
	LastSeenAt   *time.Time   `json:"last_seen_at,omitempty"`
}

func (device *Device) BeforeCreate(tx *gorm.DB) error {
	return device.BaseModel.BeforeCreate(tx)
}

func NewDevice(tenantID uuid.UUID, name, refreshToken string) *Device {
	return &Device{
		TenantID:     tenantID,
		Name:         name,
		Status:       DeviceStatusActive,
		RefreshToken: refreshToken,
	}
}

func (device *Device) IsActive() bool {
	return device.Status == DeviceStatusActive
}

type DeviceAuthStatus string

const (
	DeviceAuthStatusPending  DeviceAuthStatus = "pending"
	DeviceAuthStatusApproved DeviceAuthStatus = "approved"
	DeviceAuthStatusDenied   DeviceAuthStatus = "denied"
)

type DeviceAuthorization struct {
	BaseModel
	DeviceCode   string           `gorm:"type:varchar(64);not null;uniqueIndex" json:"-"`
	UserCode     string           `gorm:"type:varchar(16);not null;uniqueIndex" json:"user_code"`
	DeviceName   string           `gorm:"not null" json:"device_name"`
	Status       DeviceAuthStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	TenantID     *uuid.UUID       `gorm:"type:uuid" json:"tenant_id,omitempty"`
	UserID       *uuid.UUID       `gorm:"type:uuid" json:"user_id,omitempty"`
	DeviceID     *uuid.UUID       `gorm:"type:uuid" json:"device_id,omitempty"`
	ExpiresAt    time.Time        `gorm:"not null;index" json:"expires_at"`
	LastPolledAt *time.Time       `json:"-"`
	ApprovedAt   *time.Time       `json:"approved_at,omitempty"`
}

func (auth *DeviceAuthorization) BeforeCreate(tx *gorm.DB) error {
	return auth.BaseModel.BeforeCreate(tx)
}

func NewDeviceAuthorization(deviceCode, userCode, deviceName string, expiresAt time.Time) *DeviceAuthorization {
	return &DeviceAuthorization{
		DeviceCode: deviceCode,
		UserCode:   userCode,
		DeviceName: deviceName,
		Status:     DeviceAuthStatusPending,
		ExpiresAt:  expiresAt,
	}
}

func (auth *DeviceAuthorization) IsExpired() bool {
	return time.Now().After(auth.ExpiresAt)
}
