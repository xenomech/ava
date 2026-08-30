package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	BaseModel
	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Tenant   *Tenant   `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	// Unique per tenant among live rows; the partial index is built in db.Migrate.
	Name     string `gorm:"type:varchar(80);not null" json:"name"`
	Position int    `gorm:"not null;default:0" json:"position"`
}

func (room *Room) BeforeCreate(tx *gorm.DB) error {
	return room.BaseModel.BeforeCreate(tx)
}

func NewRoom(tenantID uuid.UUID, name string, position int) *Room {
	return &Room{TenantID: tenantID, Name: name, Position: position}
}
