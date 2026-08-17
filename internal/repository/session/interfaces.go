package session

import (
	"context"

	"ava/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByID(ctx context.Context, id uuid.UUID) (*model.Session, error)
	GetSessionByRID(ctx context.Context, rid string) (*model.Session, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*model.Session, error)
	UpdateSessionRID(ctx context.Context, sessionID uuid.UUID, rid string) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &sessionRepository{db: db}
}
