package adapters

import (
	"errors"
	"fmt"

	"ava/hub/external/wiz"
	"ava/hub/internal/device"
)

var ErrUnknownVendor = errors.New("adapters: unknown vendor")

func Open(spec *device.Spec) (device.Device, error) {
	if spec == nil {
		return nil, errors.New("adapters: spec is required")
	}

	switch spec.Vendor {
	case device.VendorWiz:
		return wiz.New(spec.IP, spec.Timeout), nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownVendor, spec.Vendor)
	}
}

func Vendors() []device.Vendor {
	return []device.Vendor{device.VendorWiz}
}
