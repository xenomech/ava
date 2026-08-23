package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"

	"ava/api/internal/model"
	tenantrepo "ava/api/internal/repository/tenant"
	userrepo "ava/api/internal/repository/user"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"
)

// A username and a home slug are unique keys the person registering never sees.
// Both are derived from the email local part so they stay recognisable in logs
// and URLs, then disambiguated with a counter.

const (
	identifierFallback = "home"
	identifierMaxLen   = 24
	identifierAttempts = 50
)

// identifierBase reduces an email to a slug candidate that is always valid:
// lowercase, no runs of separators, and at least three characters.
func identifierBase(email string) string {
	local, _, _ := strings.Cut(email, "@")

	base := model.Slugify(local)
	if len(base) > identifierMaxLen {
		base = strings.Trim(base[:identifierMaxLen], "-")
	}

	if len(base) < 3 {
		return identifierFallback
	}

	return base
}

// candidate returns the nth name to try: the base itself first, then base-2,
// base-3 and so on, so the common case reads as the person's own name.
func candidate(base string, n int) string {
	if n == 0 {
		return base
	}

	return base + "-" + strconv.Itoa(n+1)
}

func (s *authService) freeUsername(ctx context.Context, base string) (string, error) {
	return firstFree(base, func(name string) (bool, error) {
		_, err := s.userRepo.GetUserByUsername(ctx, name)
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return true, nil
		}

		if err != nil {
			logger.Error("auth.freeUsername", logger.Err(err))

			return false, err
		}

		return false, nil
	})
}

func (s *authService) freeSlug(ctx context.Context, base string) (string, error) {
	return firstFree(base, func(slug string) (bool, error) {
		_, err := s.tenantRepo.GetBySlug(ctx, slug)
		if serrors.Is(err, tenantrepo.ErrTenantNotFound) {
			return true, nil
		}

		if err != nil {
			logger.Error("auth.freeSlug", logger.Err(err))

			return false, err
		}

		return false, nil
	})
}

// firstFree walks candidates until available says yes. If every candidate is
// taken it falls back to a random suffix rather than failing the registration —
// the name is cosmetic, and refusing a signup over it is the worse outcome. The
// unique index stays the real guard.
func firstFree(base string, available func(string) (bool, error)) (string, error) {
	for n := range identifierAttempts {
		name := candidate(base, n)

		free, err := available(name)
		if err != nil {
			return "", err
		}

		if free {
			return name, nil
		}
	}

	suffix := make([]byte, 4)
	if _, err := rand.Read(suffix); err != nil {
		return "", err
	}

	return base + "-" + hex.EncodeToString(suffix), nil
}
