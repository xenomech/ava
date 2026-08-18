package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HubStatus string

const (
	HubStatusActive  HubStatus = "active"
	HubStatusRevoked HubStatus = "revoked"
)

type Hub struct {
	BaseModel
	TenantID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Tenant       *Tenant    `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	Name         string     `gorm:"not null" json:"name"`
	Status       HubStatus  `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	RefreshToken string     `gorm:"type:varchar(64);not null;uniqueIndex" json:"-"`
	LastSeenAt   *time.Time `json:"last_seen_at,omitempty"`
}

func (hub *Hub) BeforeCreate(tx *gorm.DB) error {
	return hub.BaseModel.BeforeCreate(tx)
}

func NewHub(tenantID uuid.UUID, name, refreshToken string) *Hub {
	return &Hub{
		TenantID:     tenantID,
		Name:         name,
		Status:       HubStatusActive,
		RefreshToken: refreshToken,
	}
}

func (hub *Hub) IsActive() bool {
	return hub.Status == HubStatusActive
}

type HubAuthStatus string

const (
	HubAuthStatusPending  HubAuthStatus = "pending"
	HubAuthStatusApproved HubAuthStatus = "approved"
	HubAuthStatusDenied   HubAuthStatus = "denied"
)

type HubAuthorization struct {
	BaseModel
	DeviceCode   string        `gorm:"type:varchar(64);not null;uniqueIndex" json:"-"`
	UserCode     string        `gorm:"type:varchar(16);not null;uniqueIndex" json:"user_code"`
	HubName      string        `gorm:"not null" json:"hub_name"`
	Status       HubAuthStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	TenantID     *uuid.UUID    `gorm:"type:uuid" json:"tenant_id,omitempty"`
	UserID       *uuid.UUID    `gorm:"type:uuid" json:"user_id,omitempty"`
	HubID        *uuid.UUID    `gorm:"type:uuid" json:"hub_id,omitempty"`
	ExpiresAt    time.Time     `gorm:"not null;index" json:"expires_at"`
	LastPolledAt *time.Time    `json:"-"`
	ApprovedAt   *time.Time    `json:"approved_at,omitempty"`
}

func (auth *HubAuthorization) BeforeCreate(tx *gorm.DB) error {
	return auth.BaseModel.BeforeCreate(tx)
}

func NewHubAuthorization(deviceCode, userCode, hubName string, expiresAt time.Time) *HubAuthorization {
	return &HubAuthorization{
		DeviceCode: deviceCode,
		UserCode:   userCode,
		HubName:    hubName,
		Status:     HubAuthStatusPending,
		ExpiresAt:  expiresAt,
	}
}

func (auth *HubAuthorization) IsExpired() bool {
	return time.Now().After(auth.ExpiresAt)
}
