package flow

import (
	"encoding/json"

	"ava/internal/model"
)

type ProfileStepData struct {
	Name  string `json:"name"`
	Phone string `json:"phone,omitempty"`
}

type WorkspaceStepData struct {
	Name string `json:"name"`
}

type InviteTeamStepData struct {
	Emails []string         `json:"emails"`
	Role   model.TenantRole `json:"role,omitempty"`
}

type RoleOption struct {
	Value model.TenantRole `json:"value"`
	Label string           `json:"label"`
}

type OnboardingMetadata struct {
	InviteRoles []RoleOption `json:"invite_roles"`
}

func decodePayload[T any](data json.RawMessage) (T, error) {
	var payload T

	if len(data) == 0 {
		return payload, nil
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return payload, ErrInvalidStepData
	}

	return payload, nil
}
