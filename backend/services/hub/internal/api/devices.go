package api

import (
	"context"
	"encoding/json"
)

type SyncDeviceItem struct {
	ExternalID string          `json:"external_id"`
	Name       string          `json:"name"`
	Kind       string          `json:"kind"`
	Status     string          `json:"status"`
	State      json.RawMessage `json:"state,omitempty"`
}

type SyncDevicesRequest struct {
	Devices []SyncDeviceItem `json:"devices"`
}

type SyncedDevice struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
}

func (c *Client) SyncDevices(ctx context.Context, devices []SyncDeviceItem) ([]SyncedDevice, error) {
	var out []SyncedDevice

	body := SyncDevicesRequest{Devices: devices}

	if err := c.do(ctx, "PUT", "/hubs/devices", body, &out); err != nil {
		return nil, err
	}

	return out, nil
}
