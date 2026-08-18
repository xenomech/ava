package hub

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"ava/hub/internal/api"
	"ava/hub/internal/config"
	"ava/hub/internal/state"
)

type Hub struct {
	cfg    *config.Config
	client *api.Client
	state  *state.State
}

func New(cfg *config.Config) *Hub {
	return &Hub{
		cfg:    cfg,
		client: api.NewClient(cfg.APIBaseURL),
	}
}

func (h *Hub) Run(ctx context.Context) error {
	loaded, err := state.Load(h.cfg.StateFile)
	if err != nil {
		return err
	}

	h.state = loaded

	tokens, err := h.authorize(ctx)
	if err != nil {
		return err
	}

	h.client.SetToken(tokens.AccessToken)

	slog.Info("HUB_PAIRED",
		slog.String("device_id", tokens.Device.ID),
		slog.String("device_name", tokens.Device.Name),
		slog.String("tenant", tokens.Tenant.Slug),
	)

	return h.heartbeatLoop(ctx, tokens)
}

func (h *Hub) authorize(ctx context.Context) (*api.DeviceTokens, error) {
	if h.state.IsPaired() {
		tokens, err := h.client.RefreshDeviceToken(ctx, h.state.RefreshToken)
		if err == nil {
			return tokens, h.persist(tokens)
		}

		code := api.CodeOf(err)
		if code != api.CodeInvalidRefreshToken && code != api.CodeDeviceRevoked {
			return nil, err
		}

		slog.Warn("HUB_PAIRING_REJECTED", slog.String("code", code))

		if err := state.Clear(h.cfg.StateFile); err != nil {
			return nil, err
		}

		h.state = &state.State{}
	}

	return h.pair(ctx)
}

func (h *Hub) pair(ctx context.Context) (*api.DeviceTokens, error) {
	code, err := h.client.RequestDeviceCode(ctx, h.cfg.DeviceName)
	if err != nil {
		return nil, err
	}

	slog.Info("HUB_ACTIVATION_REQUIRED",
		slog.String("user_code", code.UserCode),
		slog.String("verification_uri", code.VerificationURI),
		slog.Int64("expires_in", code.ExpiresIn),
	)

	interval := time.Duration(code.Interval) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}

	deadline := time.Now().Add(time.Duration(code.ExpiresIn) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}

		if time.Now().After(deadline) {
			return nil, errors.New("activation code expired before it was approved")
		}

		tokens, err := h.client.PollDeviceToken(ctx, code.DeviceCode)
		if err == nil {
			return tokens, h.persist(tokens)
		}

		switch api.CodeOf(err) {
		case api.CodeAuthorizationPending:
			continue
		case api.CodeSlowDown:
			interval += time.Second

			continue
		case api.CodeAccessDenied:
			return nil, errors.New("activation was denied")
		case api.CodeExpiredToken:
			return nil, errors.New("activation code expired before it was approved")
		default:
			return nil, err
		}
	}
}

func (h *Hub) persist(tokens *api.DeviceTokens) error {
	h.state = &state.State{
		DeviceID:     tokens.Device.ID,
		DeviceName:   tokens.Device.Name,
		TenantSlug:   tokens.Tenant.Slug,
		RefreshToken: tokens.RefreshToken,
	}

	return state.Save(h.cfg.StateFile, h.state)
}

func (h *Hub) heartbeatLoop(ctx context.Context, tokens *api.DeviceTokens) error {
	renewAt := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second / 2)

	ticker := time.NewTicker(h.cfg.HeartbeatInterval)
	defer ticker.Stop()

	for {
		if err := h.client.Heartbeat(ctx); err != nil {
			slog.Warn("HEARTBEAT_FAILED", slog.String("error", err.Error()))
		}

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		if time.Now().After(renewAt) {
			refreshed, err := h.client.RefreshDeviceToken(ctx, h.state.RefreshToken)
			if err != nil {
				slog.Warn("TOKEN_REFRESH_FAILED", slog.String("error", err.Error()))

				continue
			}

			h.client.SetToken(refreshed.AccessToken)
			renewAt = time.Now().Add(time.Duration(refreshed.ExpiresIn) * time.Second / 2)

			if err := h.persist(refreshed); err != nil {
				return err
			}

			slog.Info("TOKEN_REFRESHED")
		}
	}
}
