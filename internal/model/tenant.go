package model

import (
	"regexp"

	"gorm.io/gorm"
)

type TenantRole string

const (
	TenantRoleOwner  TenantRole = "owner"
	TenantRoleAdmin  TenantRole = "admin"
	TenantRoleMember TenantRole = "member"
)

func (role TenantRole) IsValid() bool {
	switch role {
	case TenantRoleOwner, TenantRoleAdmin, TenantRoleMember:
		return true
	default:
		return false
	}
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func IsValidSlug(slug string) bool {
	return slugPattern.MatchString(slug)
}

type Tenant struct {
	BaseModel
	Name string `gorm:"not null" json:"name"`
	Slug string `gorm:"uniqueIndex;not null" json:"slug"`
}

func (tenant *Tenant) BeforeCreate(tx *gorm.DB) error {
	return tenant.BaseModel.BeforeCreate(tx)
}

func NewTenant(name, slug string) *Tenant {
	return &Tenant{
		Name: name,
		Slug: slug,
	}
}
