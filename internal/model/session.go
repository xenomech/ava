package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	BaseModel
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	DeviceName string    `gorm:"not null" json:"device_name"`
	IPAddress  string    `gorm:"not null" json:"ip_address"`
	UserAgent  string    `gorm:"not null" json:"user_agent"`
	RID        string    `gorm:"column:rid;type:varchar(64);not null;uniqueIndex" json:"-"`
	Revoked    bool      `gorm:"default:false" json:"revoked"`
	ExpiresAt  time.Time `gorm:"not null;index" json:"expires_at"`
}

func (session *Session) BeforeCreate(tx *gorm.DB) error {
	return session.BaseModel.BeforeCreate(tx)
}

func NewSession(userID uuid.UUID, deviceName, ipAddress, userAgent, rid string, expiresAt time.Time) *Session {
	id, _ := uuid.NewV7()

	return &Session{
		BaseModel:  BaseModel{ID: id},
		UserID:     userID,
		DeviceName: deviceName,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		RID:        rid,
		Revoked:    false,
		ExpiresAt:  expiresAt,
	}
}
