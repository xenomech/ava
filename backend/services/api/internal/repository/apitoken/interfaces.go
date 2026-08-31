package apitoken

import (
	"context"
	"time"

	"ava/api/internal/model"
	"ava/api/pkg/serrors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTokenNotFound = serrors.NewCoded("api_token_not_found", "token not found")
	ErrNameTaken     = serrors.NewCoded("api_token_name_taken", "you already have a token with that name")
)

type Repository interface {
	Create(ctx context.Context, token *model.APIToken) error
	ListByUser(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.APIToken, error)
	// GetByLookup finds a token by its public half for authentication; it does not filter on tenant.
	GetByLookup(ctx context.Context, lookup string) (*model.APIToken, error)
	Revoke(ctx context.Context, tenantID, userID, tokenID uuid.UUID, at time.Time) error
	Delete(ctx context.Context, tenantID, userID, tokenID uuid.UUID) error
	TouchLastUsed(ctx context.Context, tokenID uuid.UUID, at time.Time) error
	NameExists(ctx context.Context, tenantID, userID uuid.UUID, name string) (bool, error)
}

type apiTokenRepository struct {
	db *gorm.DB
}

func NewRepository(database *gorm.DB) Repository {
	return &apiTokenRepository{db: database}
}
