package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenType string

const (
	TokenTypeEmailVerification TokenType = "email_verification"
	TokenTypePasswordReset     TokenType = "password_reset"
)

type Token struct {
	BaseModel
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Token     string     `gorm:"uniqueIndex;not null" json:"token"`
	Type      TokenType  `gorm:"type:varchar(50);not null;index" json:"type"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
}

func (token *Token) BeforeCreate(tx *gorm.DB) error {
	return token.BaseModel.BeforeCreate(tx)
}

func NewToken(userID uuid.UUID, token string, tokenType TokenType, expiresAt time.Time) *Token {
	return &Token{
		UserID:    userID,
		Token:     token,
		Type:      tokenType,
		ExpiresAt: expiresAt,
	}
}
