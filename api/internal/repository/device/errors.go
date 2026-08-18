package device

import "ava/api/pkg/serrors"

var (
	ErrDeviceNotFound        = serrors.New("device not found")
	ErrAuthorizationNotFound = serrors.New("device authorization not found")
)
