package dto

import (
	"ava/pkg/wire"

	"github.com/google/uuid"
)

type SceneTargetRequest struct {
	DeviceID uuid.UUID  `json:"device_id" validate:"required"`
	Trait    wire.Trait `json:"trait" validate:"required,max=64"`
	Value    wire.Value `json:"value"`
}

type CreateSceneRequest struct {
	Name    string               `json:"name" validate:"required,max=80"`
	Targets []SceneTargetRequest `json:"targets" validate:"required,min=1,max=100,dive"`
}

type SceneTargetResponse struct {
	DeviceID uuid.UUID  `json:"device_id"`
	Trait    wire.Trait `json:"trait"`
	Value    wire.Value `json:"value"`
}

type SceneResponse struct {
	ID       uuid.UUID             `json:"id"`
	RoomID   uuid.UUID             `json:"room_id"`
	Name     string                `json:"name"`
	Position int                   `json:"position"`
	Targets  []SceneTargetResponse `json:"targets"`
}
