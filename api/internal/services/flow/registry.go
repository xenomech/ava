package flow

import (
	"ava/api/internal/model"
	"ava/api/internal/serrors"
)

type StepDefinition struct {
	ID          string
	Title       string
	Description string
	Skippable   bool
}

type Definition struct {
	Type     string
	Steps    []StepDefinition
	Metadata any
}

var inviteRoles = []RoleOption{
	{Value: model.TenantRoleAdmin, Label: "Admin"},
	{Value: model.TenantRoleMember, Label: "Member"},
}

var registry = map[string]Definition{
	"onboarding": {
		Type: "onboarding",
		Steps: []StepDefinition{
			{
				ID:          "profile",
				Title:       "Confirm your details",
				Description: "Tell us how your name should appear to your teammates",
				Skippable:   false,
			},
			{
				ID:          "workspace",
				Title:       "Name your workspace",
				Description: "Choose the name your team will see for this workspace",
				Skippable:   true,
			},
			{
				ID:          "invite_team",
				Title:       "Invite your team",
				Description: "Invite the people you work with, or do it later",
				Skippable:   true,
			},
		},
		Metadata: OnboardingMetadata{
			InviteRoles: inviteRoles,
		},
	},
}

func GetDefinition(flowType string) (*Definition, error) {
	def, ok := registry[flowType]
	if !ok {
		return nil, serrors.New("invalid flow type: " + flowType)
	}

	return &def, nil
}
