package session

import (
	"context"
	"errors"
	"time"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *sessionRepository) CreateSession(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, tenantID, id uuid.UUID) (*model.Session, error) {
	var session model.Session

	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSessionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) GetSessionByRID(ctx context.Context, tenantID uuid.UUID, rid string) (*model.Session, error) {
	var session model.Session

	err := r.db.WithContext(ctx).Where("tenant_id = ? AND rid = ?", tenantID, rid).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSessionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) RevokeSession(ctx context.Context, tenantID, sessionID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&model.Session{}).Where("tenant_id = ? AND id = ?", tenantID, sessionID).Updates(map[string]any{
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

func (r *sessionRepository) RevokeAllUserSessions(ctx context.Context, tenantID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Session{}).Where("tenant_id = ? AND user_id = ? AND revoked = ?", tenantID, userID, false).Updates(map[string]any{
		"revoked":    true,
		"updated_at": time.Now(),
	}).Error
}

func (r *sessionRepository) RevokeAllUserSessionsGlobal(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Session{}).Where("user_id = ? AND revoked = ?", userID, false).Updates(map[string]any{
		"revoked":    true,
		"updated_at": time.Now(),
	}).Error
}

func (r *sessionRepository) GetUserSessions(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.Session, error) {
	var sessions []*model.Session

	err := r.db.WithContext(ctx).Where("tenant_id = ? AND user_id = ? AND revoked = ? AND expires_at > ?", tenantID, userID, false, time.Now()).Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *sessionRepository) UpdateSessionRID(ctx context.Context, tenantID, sessionID uuid.UUID, rid string) error {
	result := r.db.WithContext(ctx).Model(&model.Session{}).Where("tenant_id = ? AND id = ?", tenantID, sessionID).Updates(map[string]any{
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

func (r *sessionRepository) LatestTenantForUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var session model.Session

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, ErrSessionNotFound
	}

	if err != nil {
		return uuid.Nil, err
	}

	return session.TenantID, nil
}
