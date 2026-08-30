package flow

import (
	"encoding/json"
)

type HomeStepData struct {
	Name string `json:"name"`
}

// HubStepData carries nothing. Pairing happens through the hub activation
// endpoint, which the client already calls; this step only records that a hub
// arrived, so its handler reads the hub list rather than a payload.
type HubStepData struct{}

type OnboardingMetadata struct{}

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
