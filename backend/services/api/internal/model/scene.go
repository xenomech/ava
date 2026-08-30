package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Scene is a snapshot of one room, so it replays through the ordinary batch apply path as a list of targets.
type Scene struct {
	BaseModel
	TenantID uuid.UUID     `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Tenant   *Tenant       `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	RoomID   uuid.UUID     `gorm:"type:uuid;not null;index" json:"room_id"`
	Room     *Room         `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"-"`
	Name     string        `gorm:"type:varchar(80);not null" json:"name"`
	Position int           `gorm:"not null;default:0" json:"position"`
	Targets  []SceneTarget `gorm:"foreignKey:SceneID;constraint:OnDelete:CASCADE" json:"targets"`
}

// SceneTarget is one frozen trait of one device; Value is raw JSON as a trait may be bool, number or string.
type SceneTarget struct {
	BaseModel
	SceneID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"scene_id"`
	DeviceID uuid.UUID       `gorm:"type:uuid;not null;index" json:"device_id"`
	Device   *Device         `gorm:"foreignKey:DeviceID;constraint:OnDelete:CASCADE" json:"-"`
	Trait    string          `gorm:"type:varchar(64);not null" json:"trait"`
	Value    json.RawMessage `gorm:"type:jsonb;not null;default:'null'" json:"value"`
}

func (scene *Scene) BeforeCreate(tx *gorm.DB) error {
	return scene.BaseModel.BeforeCreate(tx)
}

func (target *SceneTarget) BeforeCreate(tx *gorm.DB) error {
	return target.BaseModel.BeforeCreate(tx)
}

func NewScene(tenantID, roomID uuid.UUID, name string, position int) *Scene {
	return &Scene{TenantID: tenantID, RoomID: roomID, Name: name, Position: position}
}
