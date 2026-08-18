package device

import (
	"context"
	"time"

	"ava/api/config"
	"ava/api/internal/dto"
	"ava/api/internal/model"
	devicerepo "ava/api/internal/repository/device"
	"ava/api/pkg/logger"
	"ava/api/pkg/serrors"

	"github.com/google/uuid"
)

func (s *deviceService) RequestCode(ctx context.Context, req *dto.DeviceCodeRequest) (*dto.DeviceCodeResponse, error) {
	cfg := config.GetConfig()

	deviceCode, err := generateSecret()
	if err != nil {
		logger.Error("device.RequestCode", logger.Err(err))

		return nil, err
	}

	userCode, err := generateUserCode()
	if err != nil {
		logger.Error("device.RequestCode", logger.Err(err))

		return nil, err
	}

	auth := model.NewDeviceAuthorization(deviceCode, userCode, req.DeviceName, time.Now().Add(cfg.DeviceCodeExpiry))

	if err := s.deviceRepo.CreateAuthorization(ctx, auth); err != nil {
		logger.Error("device.RequestCode", logger.Err(err))

		return nil, err
	}

	return &dto.DeviceCodeResponse{
		DeviceCode:      deviceCode,
		UserCode:        userCode,
		VerificationURI: cfg.AppURL + "/activate",
		ExpiresIn:       int64(cfg.DeviceCodeExpiry.Seconds()),
		Interval:        int64(cfg.DevicePollInterval.Seconds()),
	}, nil
}

func (s *deviceService) Poll(ctx context.Context, deviceCode string) (*dto.DeviceTokenResponse, error) {
	cfg := config.GetConfig()

	auth, err := s.deviceRepo.GetAuthorizationByDeviceCode(ctx, deviceCode)
	if err != nil {
		if serrors.Is(err, devicerepo.ErrAuthorizationNotFound) {
			return nil, ErrInvalidCode
		}

		logger.Error("device.Poll", logger.Err(err))

		return nil, err
	}

	if auth.LastPolledAt != nil && time.Since(*auth.LastPolledAt) < cfg.DevicePollInterval {
		return nil, ErrSlowDown
	}

	if err := s.deviceRepo.TouchAuthorizationPoll(ctx, auth.ID, time.Now()); err != nil {
		logger.Error("device.Poll", logger.Err(err))

		return nil, err
	}

	switch auth.Status {
	case model.DeviceAuthStatusDenied:
		return nil, ErrAccessDenied
	case model.DeviceAuthStatusPending:
		if auth.IsExpired() {
			return nil, ErrExpiredCode
		}

		return nil, ErrAuthorizationPending
	case model.DeviceAuthStatusApproved:
	}

	if auth.DeviceID == nil {
		return nil, ErrInvalidCode
	}

	device, err := s.deviceRepo.GetByIDUnscoped(ctx, *auth.DeviceID)
	if err != nil {
		logger.Error("device.Poll", logger.Err(err))

		return nil, err
	}

	return s.issueTokens(ctx, device)
}

func (s *deviceService) Refresh(ctx context.Context, refreshToken string) (*dto.DeviceTokenResponse, error) {
	device, err := s.deviceRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		if serrors.Is(err, devicerepo.ErrDeviceNotFound) {
			return nil, ErrInvalidRefreshToken
		}

		logger.Error("device.Refresh", logger.Err(err))

		return nil, err
	}

	if !device.IsActive() {
		return nil, ErrDeviceRevoked
	}

	return s.issueTokens(ctx, device)
}

func (s *deviceService) Activate(ctx context.Context, tenantID, userID uuid.UUID, userCode string) (*dto.DeviceResponse, error) {
	auth, err := s.deviceRepo.GetAuthorizationByUserCode(ctx, normalizeUserCode(userCode))
	if err != nil {
		if serrors.Is(err, devicerepo.ErrAuthorizationNotFound) {
			return nil, ErrInvalidCode
		}

		logger.Error("device.Activate", logger.Err(err))

		return nil, err
	}

	if auth.Status != model.DeviceAuthStatusPending {
		return nil, ErrInvalidCode
	}

	if auth.IsExpired() {
		return nil, ErrExpiredCode
	}

	refreshToken, err := generateSecret()
	if err != nil {
		logger.Error("device.Activate", logger.Err(err))

		return nil, err
	}

	device := model.NewDevice(tenantID, auth.DeviceName, refreshToken)

	now := time.Now()
	auth.Status = model.DeviceAuthStatusApproved
	auth.TenantID = &tenantID
	auth.UserID = &userID
	auth.ApprovedAt = &now

	if err := s.deviceRepo.ApproveWithDevice(ctx, auth, device); err != nil {
		if serrors.Is(err, devicerepo.ErrAuthorizationNotFound) {
			return nil, ErrInvalidCode
		}

		logger.Error("device.Activate", logger.Err(err))

		return nil, err
	}

	return toDeviceResponse(device), nil
}

func (s *deviceService) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.DeviceResponse, error) {
	devices, err := s.deviceRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		logger.Error("device.ListByTenant", logger.Err(err))

		return nil, err
	}

	responses := make([]*dto.DeviceResponse, 0, len(devices))
	for _, device := range devices {
		responses = append(responses, toDeviceResponse(device))
	}

	return responses, nil
}

func (s *deviceService) Revoke(ctx context.Context, tenantID, deviceID uuid.UUID) error {
	if err := s.deviceRepo.Revoke(ctx, tenantID, deviceID); err != nil {
		if serrors.Is(err, devicerepo.ErrDeviceNotFound) {
			return ErrDeviceNotFound
		}

		logger.Error("device.Revoke", logger.Err(err))

		return err
	}

	return nil
}

func (s *deviceService) Heartbeat(ctx context.Context, deviceID uuid.UUID) error {
	if err := s.deviceRepo.TouchLastSeen(ctx, deviceID, time.Now()); err != nil {
		logger.Error("device.Heartbeat", logger.Err(err))

		return err
	}

	return nil
}

func (s *deviceService) ValidateDevice(ctx context.Context, deviceID uuid.UUID) (*model.Device, error) {
	device, err := s.deviceRepo.GetByIDUnscoped(ctx, deviceID)
	if err != nil {
		if serrors.Is(err, devicerepo.ErrDeviceNotFound) {
			return nil, ErrDeviceNotFound
		}

		return nil, err
	}

	if !device.IsActive() {
		return nil, ErrDeviceRevoked
	}

	return device, nil
}

func (s *deviceService) issueTokens(ctx context.Context, device *model.Device) (*dto.DeviceTokenResponse, error) {
	accessToken, err := s.tokenManager.GenerateDeviceToken(device)
	if err != nil {
		logger.Error("device.issueTokens", logger.Err(err))

		return nil, err
	}

	refreshToken, err := generateSecret()
	if err != nil {
		logger.Error("device.issueTokens", logger.Err(err))

		return nil, err
	}

	if err := s.deviceRepo.UpdateRefreshToken(ctx, device.ID, refreshToken); err != nil {
		logger.Error("device.issueTokens", logger.Err(err))

		return nil, err
	}

	tenant, err := s.tenantRepo.GetByID(ctx, device.TenantID)
	if err != nil {
		logger.Error("device.issueTokens", logger.Err(err))

		return nil, err
	}

	return &dto.DeviceTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenManager.GetDeviceExpiry().Seconds()),
		Device:       *toDeviceResponse(device),
		Tenant: dto.TenantSummary{
			ID:   tenant.ID,
			Name: tenant.Name,
			Slug: tenant.Slug,
		},
	}, nil
}
