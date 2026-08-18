package device

import "ava/api/pkg/serrors"

var ErrDeviceNotFound = serrors.NewCoded("device_not_found", "device not found")
