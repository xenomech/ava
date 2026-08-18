package tenant

import "ava/api/pkg/serrors"

var (
	ErrTenantNotFound      = serrors.New("tenant not found")
	ErrTenantAlreadyExists = serrors.New("tenant already exists")
	ErrUserAlreadyExists   = serrors.New("user already exists")
)
