package auth

import "ava/api/pkg/serrors"

var (
	ErrInvalidCredentials      = serrors.NewCoded("invalid_credentials", "invalid credentials")
	ErrSessionExpired          = serrors.NewCoded("session_expired", "session has expired or been revoked")
	ErrUserAlreadyExists       = serrors.NewCoded("user_already_exists", "user already exists")
	ErrInvalidToken            = serrors.NewCoded("invalid_token", "invalid or expired token")
	ErrSessionRevoked          = serrors.NewCoded("session_revoked", "session has been revoked")
	ErrUserNotFound            = serrors.NewCoded("user_not_found", "user not found")
	ErrSessionNotFound         = serrors.NewCoded("session_not_found", "session not found")
	ErrAccessDenied            = serrors.NewCoded("access_denied", "access denied")
	ErrNoTenantMembership      = serrors.NewCoded("no_tenant_membership", "user does not belong to any tenant")
	ErrTenantAlreadyExists     = serrors.NewCoded("tenant_already_exists", "tenant slug is already taken")
	ErrTenantSelectionRequired = serrors.NewCoded("tenant_selection_required", "multiple tenants available, specify tenant_slug")
	ErrInvalidSlug             = serrors.NewCoded("invalid_slug", "tenant slug must be lowercase alphanumeric words separated by hyphens")
	ErrInviteInvalid           = serrors.NewCoded("invite_invalid", "invitation is invalid or has expired")
	ErrEmailNotVerified        = serrors.NewCoded("email_not_verified", "email not verified")
	ErrPasswordMismatch        = serrors.NewCoded("password_mismatch", "current password is incorrect")
)
