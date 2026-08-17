package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Email           string     `gorm:"uniqueIndex;not null" json:"email"`
	Username        string     `gorm:"uniqueIndex;not null" json:"username"`
	Name            string     `gorm:"not null" json:"name"`
	Phone           string     `json:"phone,omitempty"`
	Password        string     `gorm:"not null" json:"-"`
	EmailVerified   bool       `gorm:"default:false" json:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	return user.BaseModel.BeforeCreate(tx)
}

func NewUser(email, username, name, phone, hashedPassword string) *User {
	return &User{
		Email:         email,
		Username:      username,
		Name:          name,
		Phone:         phone,
		Password:      hashedPassword,
		EmailVerified: false,
	}
}
