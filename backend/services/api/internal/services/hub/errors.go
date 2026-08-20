package hub

import "ava/api/pkg/serrors"

var (
	ErrAuthorizationPending = serrors.NewCoded("authorization_pending", "authorization is still pending")
	ErrSlowDown             = serrors.NewCoded("slow_down", "polling too frequently")
	ErrAccessDenied         = serrors.NewCoded("access_denied", "authorization was denied")
	ErrExpiredCode          = serrors.NewCoded("expired_token", "the code has expired")
	ErrInvalidCode          = serrors.NewCoded("invalid_code", "unknown or already used code")
	ErrHubNotFound          = serrors.NewCoded("device_not_found", "hub not found")
	ErrHubRevoked           = serrors.NewCoded("hub_revoked", "hub has been revoked")
	ErrInvalidRefreshToken  = serrors.NewCoded("invalid_refresh_token", "invalid hub refresh token")
)
