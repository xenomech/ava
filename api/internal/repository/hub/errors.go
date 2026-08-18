package hub

import "ava/api/pkg/serrors"

var (
	ErrHubNotFound           = serrors.New("hub not found")
	ErrAuthorizationNotFound = serrors.New("hub authorization not found")
)
