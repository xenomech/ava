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

func TestEveryListedVendorCanBeOpened(t *testing.T) {
	for _, vendor := range adapters.Vendors() {
		if _, err := adapters.Open(&device.Spec{Vendor: vendor, IP: "192.168.1.50"}); err != nil {
			t.Errorf("%s is listed but cannot be opened: %v", vendor, err)
		}
	}
}
