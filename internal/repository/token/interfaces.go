package token

import (
	"context"

	"ava/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateToken(ctx context.Context, token *model.Token) error
	GetValidToken(ctx context.Context, tokenString string, tokenType model.TokenType) (*model.Token, error)
	MarkTokenUsed(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpiredTokens(ctx context.Context) (int64, error)
}

type tokenRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &tokenRepository{db: db}
}
