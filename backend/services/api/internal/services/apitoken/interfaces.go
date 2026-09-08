package apitoken

import (
	"context"
	"time"

	"ava/api/internal/model"
	"ava/api/pkg/serrors"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken  = serrors.NewCoded("invalid_api_token", "invalid API token")
	ErrTokenRevoked  = serrors.NewCoded("api_token_revoked", "this token has been revoked")
	ErrTokenExpired  = serrors.NewCoded("api_token_expired", "this token has expired")
	ErrUnknownScope  = serrors.NewCoded("unknown_scope", "unknown scope")
	ErrNoScopes      = serrors.NewCoded("no_scopes", "a token needs at least one scope")
	ErrAccessDenied  = serrors.NewCoded("access_denied", "you no longer have access to this home")
	ErrNameTaken     = serrors.NewCoded("api_token_name_taken", "you already have a token with that name")
	ErrTokenNotFound = serrors.NewCoded("api_token_not_found", "token not found")
)

// Authenticated is everything a request needs about the caller behind a token.
type Authenticated struct {
	Token *model.APIToken
	User  *model.User
	Role  model.TenantRole
}

type Service interface {
	Create(
		ctx context.Context,
		tenantID, userID uuid.UUID,
		name string,
		scopes []string,
		expiresAt *time.Time,
	) (*model.APIToken, string, error)
	List(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.APIToken, error)
	Revoke(ctx context.Context, tenantID, userID, tokenID uuid.UUID) error
	Delete(ctx context.Context, tenantID, userID, tokenID uuid.UUID) error
	// Authenticate resolves a presented token, or fails without saying which part was wrong.
	Authenticate(ctx context.Context, presented string) (*Authenticated, error)
}
