package auth

import (
	"context"
	"time"

	"ava/internal/dto"
	"ava/internal/model"
	tenantrepo "ava/internal/repository/tenant"
	"ava/internal/services/auth/jwt"
	"ava/pkg/logger"
	"ava/pkg/serrors"

	"github.com/google/uuid"
)

func (s *authService) activeMemberships(ctx context.Context, userID uuid.UUID) ([]*model.TenantMembership, error) {
	memberships, err := s.membershipRepo.ListByUser(ctx, userID)
	if err != nil {
		logger.Error("failed to list user memberships", logger.Err(err))

		return nil, err
	}

	active := make([]*model.TenantMembership, 0, len(memberships))

	for _, membership := range memberships {
		if membership.IsActive() {
			active = append(active, membership)
		}
	}

	return active, nil
}

func (s *authService) selectMembership(ctx context.Context, memberships []*model.TenantMembership, tenantSlug string) (*model.TenantMembership, error) {
	if tenantSlug != "" {
		tenant, err := s.tenantRepo.GetBySlug(ctx, tenantSlug)
		if err != nil {
			if serrors.Is(err, tenantrepo.ErrTenantNotFound) {
				return nil, ErrAccessDenied
			}

			logger.Error("failed to get tenant by slug", logger.Err(err))

			return nil, err
		}

		for _, membership := range memberships {
			if membership.TenantID == tenant.ID {
				return membership, nil
			}
		}

		return nil, ErrAccessDenied
	}

	if len(memberships) == 1 {
		return memberships[0], nil
	}

	return nil, ErrTenantSelectionRequired
}

func (s *authService) tenantSummaries(ctx context.Context, memberships []*model.TenantMembership) ([]dto.TenantSummary, error) {
	ids := make([]uuid.UUID, 0, len(memberships))
	roles := make(map[uuid.UUID]model.TenantRole, len(memberships))

	for _, membership := range memberships {
		ids = append(ids, membership.TenantID)
		roles[membership.TenantID] = membership.Role
	}

	tenants, err := s.tenantRepo.ListByIDs(ctx, ids)
	if err != nil {
		logger.Error("failed to list tenants", logger.Err(err))

		return nil, err
	}

	summaries := make([]dto.TenantSummary, 0, len(tenants))

	for _, tenant := range tenants {
		summaries = append(summaries, dto.TenantSummary{
			ID:   tenant.ID,
			Name: tenant.Name,
			Slug: tenant.Slug,
			Role: roles[tenant.ID],
		})
	}

	return summaries, nil
}

func (s *authService) issueSession(ctx context.Context, user *model.User, membership *model.TenantMembership, deviceInfo dto.DeviceInfo) (*dto.AuthResponse, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, membership.TenantID)
	if err != nil {
		logger.Error("failed to get tenant by ID", logger.Err(err))

		return nil, err
	}

	rid := uuid.NewString()

	session := model.NewSession(
		tenant.ID,
		user.ID,
		deviceInfo.DeviceName,
		deviceInfo.IPAddress,
		deviceInfo.UserAgent,
		rid,
		time.Now().Add(s.tokenManager.GetRefreshExpiry()),
	)

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		logger.Error("failed to create session", logger.Err(err))

		return nil, err
	}

	tokens, err := s.issueTokens(user, membership, session.ID, rid)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: *s.userToResponse(user),
		Tenant: &dto.TenantSummary{
			ID:   tenant.ID,
			Name: tenant.Name,
			Slug: tenant.Slug,
			Role: membership.Role,
		},
		Tokens: tokens,
	}, nil
}

func (s *authService) issueTokens(user *model.User, membership *model.TenantMembership, sessionID uuid.UUID, rid string) (*dto.TokenResponse, error) {
	accessToken, err := s.tokenManager.GenerateToken(user, membership.TenantID, membership.Role, sessionID, jwt.AccessToken, "")
	if err != nil {
		logger.Error("failed to generate access token", logger.Err(err))

		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateToken(user, membership.TenantID, membership.Role, sessionID, jwt.RefreshToken, rid)
	if err != nil {
		logger.Error("failed to generate refresh token", logger.Err(err))

		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenManager.GetAccessExpiry().Seconds()),
	}, nil
}

func (s *authService) userToResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		Name:            user.Name,
		Phone:           user.Phone,
		EmailVerified:   user.EmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:       user.CreatedAt,
	}
}
