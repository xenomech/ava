package adapters

import (
	"errors"
	"fmt"

	"ava/hub/external/tuya"
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
	case device.VendorTuya:
		return tuya.New(&tuya.Config{
			ID:           spec.ID,
			IP:           spec.IP,
			Name:         spec.Name,
			LocalKey:     spec.LocalKey,
			Capabilities: spec.Capabilities,
			Timeout:      spec.Timeout,
		})
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownVendor, spec.Vendor)
	}
}

func Vendors() []device.Vendor {
	return []device.Vendor{device.VendorWiz, device.VendorTuya}
}
