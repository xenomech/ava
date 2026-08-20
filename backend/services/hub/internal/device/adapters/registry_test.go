package adapters_test

import (
	"errors"
	"testing"

	"ava/hub/internal/device"
	"ava/hub/internal/device/adapters"
)

func TestOpenReturnsTheRightAdapterPerVendor(t *testing.T) {
	tests := []struct {
		name   string
		spec   device.Spec
		vendor string
	}{
		{
			name:   "wiz",
			spec:   device.Spec{Vendor: device.VendorWiz, IP: "192.168.1.50"},
			vendor: "wiz",
		},
		{
			name:   "tuya",
			spec:   device.Spec{Vendor: device.VendorTuya, ID: "bf01", IP: "192.168.1.60", LocalKey: "0123456789abcdef"},
			vendor: "tuya",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dev, err := adapters.Open(&tc.spec)
			if err != nil {
				t.Fatalf("Open: %v", err)
			}

			if got := dev.Info().Vendor; got != tc.vendor {
				t.Errorf("vendor = %s, want %s", got, tc.vendor)
			}

			if got := dev.Info().IP; got != tc.spec.IP {
				t.Errorf("ip = %s", got)
			}
		})
	}
}

func TestOpenRejectsAnUnknownVendor(t *testing.T) {
	_, err := adapters.Open(&device.Spec{Vendor: "hue"})

	if !errors.Is(err, adapters.ErrUnknownVendor) {
		t.Errorf("got %v", err)
	}
}

func TestOpenRejectsANilSpec(t *testing.T) {
	if _, err := adapters.Open(nil); err == nil {
		t.Error("expected an error")
	}
}

func TestOpenPropagatesAdapterValidation(t *testing.T) {
	_, err := adapters.Open(&device.Spec{Vendor: device.VendorTuya, ID: "bf01", LocalKey: "short"})

	if err == nil {
		t.Fatal("a bad local key must not produce a device")
	}
}

func TestCapabilitiesFlowThroughTheSpec(t *testing.T) {
	dev, err := adapters.Open(&device.Spec{
		Vendor:       device.VendorTuya,
		ID:           "bf01",
		LocalKey:     "0123456789abcdef",
		Capabilities: device.CapabilityBrightness,
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	if !dev.Capabilities().Has(device.CapabilityBrightness) {
		t.Error("capabilities did not reach the adapter")
	}
}
