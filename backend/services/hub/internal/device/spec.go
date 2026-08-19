package device

import "time"

type Vendor string

const (
	VendorWiz  Vendor = "wiz"
	VendorTuya Vendor = "tuya"
)

type Spec struct {
	Vendor       Vendor
	ID           string
	Name         string
	IP           string
	LocalKey     string
	Capabilities Capability
	Timeout      time.Duration
}
