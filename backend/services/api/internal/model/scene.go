package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Scene is a remembered arrangement of one room: what each device was doing at
// the moment it was saved.
//
// It is a snapshot rather than a rule. Nothing here describes intent — no "warm
// in the evening", no conditions — because a scene is captured from a room the
// person has already set up by hand, and the only honest thing to store is what
// they had. That also means a scene can be replayed by the ordinary batch apply
// path with no special casing: it is already a list of targets.
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

// SceneTarget is one trait of one device, frozen.
//
// Value is raw JSON rather than a typed column because a trait can be a
// boolean, a number or a string, and wire.Value already knows how to read each
// of those back out.
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
