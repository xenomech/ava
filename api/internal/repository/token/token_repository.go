package token

import (
	"context"
	"errors"
	"time"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *tokenRepository) CreateToken(ctx context.Context, token *model.Token) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) GetValidToken(ctx context.Context, tokenString string, tokenType model.TokenType) (*model.Token, error) {
	var token model.Token

	err := r.db.WithContext(ctx).
		Where("token = ? AND type = ? AND used_at IS NULL AND expires_at > ?", tokenString, tokenType, time.Now()).
		First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTokenNotFound
	}

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *tokenRepository) MarkTokenUsed(ctx context.Context, tokenID uuid.UUID) error {
	now := time.Now()

	result := r.db.WithContext(ctx).Model(&model.Token{}).Where("id = ?", tokenID).Updates(map[string]any{
		"used_at":    &now,
		"updated_at": now,
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTokenNotFound
	}

	return nil
}

func (r *tokenRepository) DeleteExpiredTokens(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&model.Token{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
