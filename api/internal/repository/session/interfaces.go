package session

import (
	"context"

	"ava/api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByID(ctx context.Context, tenantID, id uuid.UUID) (*model.Session, error)
	GetSessionByRID(ctx context.Context, tenantID uuid.UUID, rid string) (*model.Session, error)
	RevokeSession(ctx context.Context, tenantID, sessionID uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, tenantID, userID uuid.UUID) error
	RevokeAllUserSessionsGlobal(ctx context.Context, userID uuid.UUID) error
	GetUserSessions(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.Session, error)
	LatestTenantForUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	UpdateSessionRID(ctx context.Context, tenantID, sessionID uuid.UUID, rid string) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &sessionRepository{db: db}
}
