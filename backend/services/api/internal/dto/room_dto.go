package dto

import "github.com/google/uuid"

type CreateRoomRequest struct {
	Name string `json:"name" validate:"required,max=80"`
}

type UpdateRoomRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,max=80"`
	Position *int    `json:"position,omitempty" validate:"omitempty,min=0"`
}

type RoomResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Position int       `json:"position"`
}
