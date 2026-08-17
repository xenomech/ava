package membership

import "ava/pkg/serrors"

var (
	ErrMembershipNotFound      = serrors.New("membership not found")
	ErrMembershipAlreadyExists = serrors.New("membership already exists")
)
