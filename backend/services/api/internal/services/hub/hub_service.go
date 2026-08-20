package hub

import (
	"context"
	"time"

	"ava/api/config"
	"ava/api/internal/dto"
	"ava/api/internal/model"
	hubrepo "ava/api/internal/repository/hub"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"

	"github.com/google/uuid"
)

func (s *hubService) RequestCode(ctx context.Context, req *dto.DeviceCodeRequest) (*dto.DeviceCodeResponse, error) {
	cfg := config.GetConfig()

	deviceCode, err := generateSecret()
	if err != nil {
		logger.Error("hub.RequestCode", logger.Err(err))

		return nil, err
	}

	userCode, err := generateUserCode()
	if err != nil {
		logger.Error("hub.RequestCode", logger.Err(err))

		return nil, err
	}

	auth := model.NewHubAuthorization(deviceCode, userCode, req.HubName, time.Now().Add(cfg.DeviceCodeExpiry))

	if err := s.hubRepo.CreateAuthorization(ctx, auth); err != nil {
		logger.Error("hub.RequestCode", logger.Err(err))

		return nil, err
	}

	return &dto.DeviceCodeResponse{
		DeviceCode:      deviceCode,
		UserCode:        userCode,
		VerificationURI: cfg.AppURL + "/activate",
		ExpiresIn:       int64(cfg.DeviceCodeExpiry.Seconds()),
		Interval:        int64(cfg.HubPollInterval.Seconds()),
	}, nil
}

func (s *hubService) Poll(ctx context.Context, deviceCode string) (*dto.HubTokenResponse, error) {
	cfg := config.GetConfig()

	auth, err := s.hubRepo.GetAuthorizationByDeviceCode(ctx, deviceCode)
	if err != nil {
		if serrors.Is(err, hubrepo.ErrAuthorizationNotFound) {
			return nil, ErrInvalidCode
		}

		logger.Error("hub.Poll", logger.Err(err))

		return nil, err
	}

	if auth.LastPolledAt != nil && time.Since(*auth.LastPolledAt) < cfg.HubPollInterval {
		return nil, ErrSlowDown
	}

	if err := s.hubRepo.TouchAuthorizationPoll(ctx, auth.ID, time.Now()); err != nil {
		logger.Error("hub.Poll", logger.Err(err))

		return nil, err
	}

	switch auth.Status {
	case model.HubAuthStatusDenied:
		return nil, ErrAccessDenied
	case model.HubAuthStatusPending:
		if auth.IsExpired() {
			return nil, ErrExpiredCode
		}

		return nil, ErrAuthorizationPending
	case model.HubAuthStatusApproved:
	}

	if auth.HubID == nil {
		return nil, ErrInvalidCode
	}

	hub, err := s.hubRepo.GetByIDUnscoped(ctx, *auth.HubID)
	if err != nil {
		logger.Error("hub.Poll", logger.Err(err))

		return nil, err
	}

	return s.issueTokens(ctx, hub)
}

func (s *hubService) Refresh(ctx context.Context, refreshToken string) (*dto.HubTokenResponse, error) {
	hub, err := s.hubRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		if serrors.Is(err, hubrepo.ErrHubNotFound) {
			return nil, ErrInvalidRefreshToken
		}

		logger.Error("hub.Refresh", logger.Err(err))

		return nil, err
	}

	if !hub.IsActive() {
		return nil, ErrHubRevoked
	}

	return s.issueTokens(ctx, hub)
}

func (s *hubService) Activate(ctx context.Context, tenantID, userID uuid.UUID, userCode string) (*dto.HubResponse, error) {
	auth, err := s.hubRepo.GetAuthorizationByUserCode(ctx, normalizeUserCode(userCode))
	if err != nil {
		if serrors.Is(err, hubrepo.ErrAuthorizationNotFound) {
			return nil, ErrInvalidCode
		}

		logger.Error("hub.Activate", logger.Err(err))

		return nil, err
	}

	if auth.Status != model.HubAuthStatusPending {
		return nil, ErrInvalidCode
	}

	if auth.IsExpired() {
		return nil, ErrExpiredCode
	}

	refreshToken, err := generateSecret()
	if err != nil {
		logger.Error("hub.Activate", logger.Err(err))

		return nil, err
	}

	hub := model.NewHub(tenantID, auth.HubName, refreshToken)

	now := time.Now()
	auth.Status = model.HubAuthStatusApproved
	auth.TenantID = &tenantID
	auth.UserID = &userID
	auth.ApprovedAt = &now

	if err := s.hubRepo.ApproveWithDevice(ctx, auth, hub); err != nil {
		if serrors.Is(err, hubrepo.ErrAuthorizationNotFound) {
			return nil, ErrInvalidCode
		}

		logger.Error("hub.Activate", logger.Err(err))

		return nil, err
	}

	return toHubResponse(hub), nil
}

func (s *hubService) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.HubResponse, error) {
	hubs, err := s.hubRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		logger.Error("hub.ListByTenant", logger.Err(err))

		return nil, err
	}

	responses := make([]*dto.HubResponse, 0, len(hubs))
	for _, hub := range hubs {
		responses = append(responses, toHubResponse(hub))
	}

	return responses, nil
}

func (s *hubService) Revoke(ctx context.Context, tenantID, hubID uuid.UUID) error {
	if err := s.hubRepo.Revoke(ctx, tenantID, hubID); err != nil {
		if serrors.Is(err, hubrepo.ErrHubNotFound) {
			return ErrHubNotFound
		}

		logger.Error("hub.Revoke", logger.Err(err))

		return err
	}

	return nil
}

func (s *hubService) Heartbeat(ctx context.Context, hubID uuid.UUID) error {
	if err := s.hubRepo.TouchLastSeen(ctx, hubID, time.Now()); err != nil {
		logger.Error("hub.Heartbeat", logger.Err(err))

		return err
	}

	return nil
}

func (s *hubService) ValidateDevice(ctx context.Context, hubID uuid.UUID) (*model.Hub, error) {
	hub, err := s.hubRepo.GetByIDUnscoped(ctx, hubID)
	if err != nil {
		if serrors.Is(err, hubrepo.ErrHubNotFound) {
			return nil, ErrHubNotFound
		}

		return nil, err
	}

	if !hub.IsActive() {
		return nil, ErrHubRevoked
	}

	return hub, nil
}

func (s *hubService) issueTokens(ctx context.Context, hub *model.Hub) (*dto.HubTokenResponse, error) {
	accessToken, err := s.tokenManager.GenerateHubToken(hub)
	if err != nil {
		logger.Error("hub.issueTokens", logger.Err(err))

		return nil, err
	}

	refreshToken, err := generateSecret()
	if err != nil {
		logger.Error("hub.issueTokens", logger.Err(err))

		return nil, err
	}

	if err := s.hubRepo.UpdateRefreshToken(ctx, hub.ID, refreshToken); err != nil {
		logger.Error("hub.issueTokens", logger.Err(err))

		return nil, err
	}

	tenant, err := s.tenantRepo.GetByID(ctx, hub.TenantID)
	if err != nil {
		logger.Error("hub.issueTokens", logger.Err(err))

		return nil, err
	}

	return &dto.HubTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenManager.GetHubExpiry().Seconds()),
		Hub:          *toHubResponse(hub),
		Tenant: dto.TenantSummary{
			ID:   tenant.ID,
			Name: tenant.Name,
			Slug: tenant.Slug,
		},
	}, nil
}
