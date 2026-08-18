package dto

import (
	"encoding/json"

	"ava/api/internal/model"
)

type SubmitStepRequest struct {
	Data json.RawMessage `json:"data" validate:"required"`
}

type FlowStepResponse struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Skippable   bool             `json:"skippable"`
	Data        json.RawMessage  `json:"data"`
	Errors      model.StepErrors `json:"errors"`
}

type FlowStateResponse struct {
	FlowType    string             `json:"flow_type"`
	Status      string             `json:"status"`
	CurrentStep string             `json:"current_step"`
	Steps       []FlowStepResponse `json:"steps"`
	Metadata    json.RawMessage    `json:"metadata"`
}
