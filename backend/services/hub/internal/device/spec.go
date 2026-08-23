package device

import (
	"time"

	"ava/pkg/wire"
)

type Vendor string

const (
	VendorWiz Vendor = "wiz"
)

type Spec struct {
	Vendor       Vendor
	ID           string
	Name         string
	IP           string
	Capabilities wire.Capabilities
	Timeout      time.Duration
}
