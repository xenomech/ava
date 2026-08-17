package dto

import (
	"time"

	"ava/internal/model"

	"github.com/google/uuid"
)

type CreateTenantRequest struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required,min=3,max=40"`
}

type UpdateTenantRequest struct {
	Name string `json:"name" validate:"required"`
}

type InviteMemberRequest struct {
	Email string           `json:"email" validate:"required,email"`
	Role  model.TenantRole `json:"role" validate:"required"`
}

type UpdateMemberRoleRequest struct {
	Role model.TenantRole `json:"role" validate:"required"`
}

type TenantResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type InviteResponse struct {
	Member      MemberResponse `json:"member"`
	InviteToken string         `json:"invite_token"`
}

type MemberResponse struct {
	UserID   uuid.UUID              `json:"user_id"`
	Email    string                 `json:"email"`
	Name     string                 `json:"name"`
	Role     model.TenantRole       `json:"role"`
	Status   model.MembershipStatus `json:"status"`
	JoinedAt *time.Time             `json:"joined_at,omitempty"`
}
