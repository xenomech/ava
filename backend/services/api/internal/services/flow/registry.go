package flow

import (
	"ava/api/pkg/serrors"
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

// Onboarding asks only for the two things the product cannot work without: a home name and a hub.
var registry = map[string]Definition{
	"onboarding": {
		Type: "onboarding",
		Steps: []StepDefinition{
			{
				ID:          "home",
				Title:       "Name your home",
				Description: "This is what you will see at the top of every screen",
				Skippable:   false,
			},
			{
				ID:          "hub",
				Title:       "Connect your hub",
				Description: "The hub is the box on your network that finds and drives your lights",
				Skippable:   true,
			},
		},
		Metadata: OnboardingMetadata{},
	},
}

func GetDefinition(flowType string) (*Definition, error) {
	def, ok := registry[flowType]
	if !ok {
		return nil, serrors.New("invalid flow type: " + flowType)
	}

	return &def, nil
}
