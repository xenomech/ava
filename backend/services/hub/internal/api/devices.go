package api

import (
	"context"

	"ava/pkg/wire"
)

func (c *Client) SyncDevices(ctx context.Context, devices []wire.DeviceReport) ([]wire.SyncedDevice, error) {
	var out []wire.SyncedDevice

	body := wire.SyncRequest{Devices: devices}

	if err := c.do(ctx, "PUT", "/hubs/devices", body, &out); err != nil {
		return nil, err
	}

	return out, nil
}
