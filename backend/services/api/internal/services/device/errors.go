package device

import "ava/api/pkg/serrors"

var (
	ErrDeviceNotFound  = serrors.NewCoded("device_not_found", "device not found")
	ErrNothingToUpdate = serrors.NewCoded("nothing_to_update", "provide a name or a room")
)
