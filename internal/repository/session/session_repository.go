package session

import (
	"context"
	"errors"
	"time"

	"ava/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *sessionRepository) CreateSession(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*model.Session, error) {
	var session model.Session

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSessionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) GetSessionByRID(ctx context.Context, rid string) (*model.Session, error) {
	var session model.Session

	err := r.db.WithContext(ctx).Where("rid = ?", rid).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSessionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&model.Session{}).Where("id = ?", sessionID).Updates(map[string]any{
		"revoked":    true,
		"updated_at": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrSessionNotFound
	}

	return nil
}

func (r *sessionRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ? AND revoked = ?", userID, false).Updates(map[string]any{
		"revoked":    true,
		"updated_at": time.Now(),
	}).Error
}

func (r *sessionRepository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*model.Session, error) {
	var sessions []*model.Session

	err := r.db.WithContext(ctx).Where("user_id = ? AND revoked = ? AND expires_at > ?", userID, false, time.Now()).Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *sessionRepository) UpdateSessionRID(ctx context.Context, sessionID uuid.UUID, rid string) error {
	result := r.db.WithContext(ctx).Model(&model.Session{}).Where("id = ?", sessionID).Updates(map[string]any{
		"rid":        rid,
		"updated_at": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrSessionNotFound
	}

	return nil
}
