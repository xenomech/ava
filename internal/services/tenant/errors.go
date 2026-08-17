package tenant

import "ava/pkg/serrors"

var (
	ErrTenantNotFound      = serrors.NewCoded("tenant_not_found", "tenant not found")
	ErrTenantAlreadyExists = serrors.NewCoded("tenant_already_exists", "tenant slug is already taken")
	ErrInvalidSlug         = serrors.NewCoded("invalid_slug", "tenant slug must be lowercase alphanumeric words separated by hyphens")
	ErrInvalidRole         = serrors.NewCoded("invalid_role", "invalid tenant role")
	ErrUserNotFound        = serrors.NewCoded("user_not_found", "no account exists for that email address")
	ErrMemberNotFound      = serrors.NewCoded("member_not_found", "member not found in this tenant")
	ErrAlreadyMember       = serrors.NewCoded("already_member", "user is already a member of this tenant")
	ErrAlreadyInvited      = serrors.NewCoded("already_invited", "user has already been invited to this tenant")
	ErrLastOwner           = serrors.NewCoded("last_owner", "the last owner of a tenant cannot be removed or demoted")
)
