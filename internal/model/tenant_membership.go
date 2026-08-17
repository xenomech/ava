package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembershipStatus string

const (
	MembershipStatusInvited MembershipStatus = "invited"
	MembershipStatusActive  MembershipStatus = "active"
)

type TenantMembership struct {
	BaseModel
	TenantID        uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex:idx_membership_tenant_user" json:"tenant_id"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex:idx_membership_tenant_user;index:idx_membership_user" json:"user_id"`
	Role            TenantRole       `gorm:"type:varchar(20);not null" json:"role"`
	Status          MembershipStatus `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	InviteToken     string           `gorm:"type:varchar(64);index" json:"-"`
	InviteExpiresAt *time.Time       `json:"-"`
	InvitedByID     *uuid.UUID       `gorm:"type:uuid" json:"invited_by_id,omitempty"`
	JoinedAt        *time.Time       `json:"joined_at,omitempty"`
}

func (membership *TenantMembership) BeforeCreate(tx *gorm.DB) error {
	return membership.BaseModel.BeforeCreate(tx)
}

func NewTenantMembership(tenantID, userID uuid.UUID, role TenantRole) *TenantMembership {
	now := time.Now()

	return &TenantMembership{
		TenantID: tenantID,
		UserID:   userID,
		Role:     role,
		Status:   MembershipStatusActive,
		JoinedAt: &now,
	}
}

func NewTenantInvite(tenantID, userID, invitedByID uuid.UUID, role TenantRole, inviteToken string, expiresAt time.Time) *TenantMembership {
	return &TenantMembership{
		TenantID:        tenantID,
		UserID:          userID,
		Role:            role,
		Status:          MembershipStatusInvited,
		InviteToken:     inviteToken,
		InviteExpiresAt: &expiresAt,
		InvitedByID:     &invitedByID,
	}
}

func (membership *TenantMembership) IsActive() bool {
	return membership.Status == MembershipStatusActive
}
