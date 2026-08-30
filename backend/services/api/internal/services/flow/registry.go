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

// Onboarding asks for the two things the product cannot work without: what this
// home is called, and a hub to reach the lights through. The person's own name
// and their home's slug are already settled at registration, so re-asking them
// here would be the same question twice. Inviting people lives in Settings —
// it is not on the path to a working light.
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
