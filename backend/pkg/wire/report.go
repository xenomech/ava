package wire

type DeviceReport struct {
	ExternalID   string       `json:"external_id" validate:"required,max=128"`
	Name         string       `json:"name" validate:"required,max=100"`
	Kind         string       `json:"kind" validate:"required,max=40"`
	Status       string       `json:"status" validate:"required,oneof=online offline"`
	Vendor       string       `json:"vendor,omitempty" validate:"omitempty,max=40"`
	Model        string       `json:"model,omitempty" validate:"omitempty,max=80"`
	IP           string       `json:"ip,omitempty" validate:"omitempty,max=45"`
	Parent       string       `json:"parent,omitempty" validate:"omitempty,max=128"`
	Capabilities Capabilities `json:"capabilities"`
	State        State        `json:"state,omitempty"`
}

type SyncRequest struct {
	Devices []DeviceReport `json:"devices"`
}

type SyncedDevice struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
}
