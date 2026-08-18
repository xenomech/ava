package membership

import "ava/api/internal/serrors"

var (
	ErrMembershipNotFound      = serrors.New("membership not found")
	ErrMembershipAlreadyExists = serrors.New("membership already exists")
)
