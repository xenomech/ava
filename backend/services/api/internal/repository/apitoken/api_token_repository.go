package apitoken

import (
	"context"
	"errors"
	"time"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *apiTokenRepository) Create(ctx context.Context, token *model.APIToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *apiTokenRepository) ListByUser(
	ctx context.Context,
	tenantID, userID uuid.UUID,
) ([]*model.APIToken, error) {
	var tokens []*model.APIToken

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Order("created_at DESC").
		Find(&tokens).Error

	return tokens, err
}

func (r *apiTokenRepository) GetByLookup(ctx context.Context, lookup string) (*model.APIToken, error) {
	var token model.APIToken

	err := r.db.WithContext(ctx).Where("lookup = ?", lookup).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}

		return nil, err
	}

	return &token, nil
}

// Revoke stamps the token rather than deleting it, so a token that was used stays auditable.
func (r *apiTokenRepository) Revoke(
	ctx context.Context,
	tenantID, userID, tokenID uuid.UUID,
	at time.Time,
) error {
	result := r.db.WithContext(ctx).
		Model(&model.APIToken{}).
		Where("tenant_id = ? AND user_id = ? AND id = ? AND revoked_at IS NULL", tenantID, userID, tokenID).
		Update("revoked_at", at)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTokenNotFound
	}

	return nil
}

func (r *apiTokenRepository) Delete(ctx context.Context, tenantID, userID, tokenID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ? AND id = ?", tenantID, userID, tokenID).
		Delete(&model.APIToken{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTokenNotFound
	}

	return nil
}

// TouchLastUsed is best effort: a failure here must not fail the request the token authenticated.
func (r *apiTokenRepository) TouchLastUsed(ctx context.Context, tokenID uuid.UUID, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.APIToken{}).
		Where("id = ?", tokenID).
		Update("last_used_at", at).Error
}

func (r *apiTokenRepository) NameExists(
	ctx context.Context,
	tenantID, userID uuid.UUID,
	name string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.APIToken{}).
		Where("tenant_id = ? AND user_id = ? AND lower(name) = lower(?)", tenantID, userID, name).
		Count(&count).Error

	return count > 0, err
}
