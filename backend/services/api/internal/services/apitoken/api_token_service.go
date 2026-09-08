package apitoken

import (
	"context"
	"strings"
	"time"

	"ava/api/internal/model"
	apitokenrepo "ava/api/internal/repository/apitoken"
	membershiprepo "ava/api/internal/repository/membership"
	userrepo "ava/api/internal/repository/user"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"

	"github.com/google/uuid"
)

type apiTokenService struct {
	tokens      apitokenrepo.Repository
	users       userrepo.Repository
	memberships membershiprepo.Repository
}

func NewService(
	tokens apitokenrepo.Repository,
	users userrepo.Repository,
	memberships membershiprepo.Repository,
) Service {
	return &apiTokenService{tokens: tokens, users: users, memberships: memberships}
}

func (s *apiTokenService) Create(
	ctx context.Context,
	tenantID, userID uuid.UUID,
	name string,
	scopes []string,
	expiresAt *time.Time,
) (*model.APIToken, string, error) {
	name = strings.TrimSpace(name)

	granted, err := parseScopes(scopes)
	if err != nil {
		return nil, "", err
	}

	taken, err := s.tokens.NameExists(ctx, tenantID, userID, name)
	if err != nil {
		logger.Error("failed to check token name", logger.Err(err))

		return nil, "", err
	}

	if taken {
		return nil, "", ErrNameTaken
	}

	token, plaintext, err := model.NewAPIToken(tenantID, userID, name, granted, expiresAt)
	if err != nil {
		logger.Error("failed to mint token", logger.Err(err))

		return nil, "", err
	}

	if err := s.tokens.Create(ctx, token); err != nil {
		logger.Error("failed to store token", logger.Err(err))

		return nil, "", err
	}

	return token, plaintext, nil
}

func (s *apiTokenService) List(
	ctx context.Context,
	tenantID, userID uuid.UUID,
) ([]*model.APIToken, error) {
	return s.tokens.ListByUser(ctx, tenantID, userID)
}

func (s *apiTokenService) Revoke(ctx context.Context, tenantID, userID, tokenID uuid.UUID) error {
	err := s.tokens.Revoke(ctx, tenantID, userID, tokenID, time.Now())
	if serrors.Is(err, apitokenrepo.ErrTokenNotFound) {
		return ErrTokenNotFound
	}

	return err
}

func (s *apiTokenService) Delete(ctx context.Context, tenantID, userID, tokenID uuid.UUID) error {
	err := s.tokens.Delete(ctx, tenantID, userID, tokenID)
	if serrors.Is(err, apitokenrepo.ErrTokenNotFound) {
		return ErrTokenNotFound
	}

	return err
}

// Authenticate resolves a token; every failure returns one error so a probe cannot map valid tokens.
func (s *apiTokenService) Authenticate(ctx context.Context, presented string) (*Authenticated, error) {
	lookup, secret, ok := model.SplitToken(presented)
	if !ok {
		return nil, ErrInvalidToken
	}

	token, err := s.tokens.GetByLookup(ctx, lookup)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.MatchesSecret(secret) {
		return nil, ErrInvalidToken
	}

	if !token.Active() {
		return nil, ErrInvalidToken
	}

	// The token carries the member's role, so losing membership immediately disarms every token.
	membership, err := s.memberships.GetByTenantAndUser(ctx, token.TenantID, token.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !membership.IsActive() {
		return nil, ErrAccessDenied
	}

	user, err := s.users.GetUserByID(ctx, token.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if err := s.tokens.TouchLastUsed(ctx, token.ID, time.Now()); err != nil {
		logger.Warn("failed to record token use", logger.Err(err))
	}

	return &Authenticated{Token: token, User: user, Role: membership.Role}, nil
}

func parseScopes(requested []string) ([]model.Scope, error) {
	if len(requested) == 0 {
		return nil, ErrNoScopes
	}

	seen := make(map[string]bool, len(requested))
	granted := make([]model.Scope, 0, len(requested))

	for _, candidate := range requested {
		candidate = strings.TrimSpace(candidate)

		if !model.ValidScope(candidate) {
			return nil, ErrUnknownScope
		}

		if seen[candidate] {
			continue
		}

		seen[candidate] = true

		granted = append(granted, model.Scope(candidate))
	}

	return granted, nil
}
