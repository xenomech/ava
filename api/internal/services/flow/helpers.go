package flow

import (
	"encoding/json"
	"slices"

	"ava/api/internal/dto"
	"ava/api/internal/model"
)

func emptyObjectIfBlank(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return json.RawMessage("{}")
	}

	return raw
}

func findStep(steps []model.FlowStep, stepID string) *model.FlowStep {
	for i := range steps {
		if steps[i].StepID == stepID {
			return &steps[i]
		}
	}

	return nil
}

func findPreviousStep(steps []model.FlowStep, currentOrder int) *model.FlowStep {
	var previous *model.FlowStep

	for i := range steps {
		if steps[i].StepOrder >= currentOrder {
			continue
		}

		if previous == nil || steps[i].StepOrder > previous.StepOrder {
			previous = &steps[i]
		}
	}

	return previous
}

func findNextPendingStep(steps []model.FlowStep) *model.FlowStep {
	ordered := make([]*model.FlowStep, 0, len(steps))
	for i := range steps {
		ordered = append(ordered, &steps[i])
	}

	slices.SortFunc(ordered, func(a, b *model.FlowStep) int { return a.StepOrder - b.StepOrder })

	for _, step := range ordered {
		if step.Status == model.FlowStepStatusPending {
			return step
		}
	}

	return nil
}

func toFlowStateResponse(flow *model.Flow) *dto.FlowStateResponse {
	steps := make([]dto.FlowStepResponse, len(flow.Steps))

	for i := range flow.Steps {
		errs := flow.Steps[i].Errors
		if errs == nil {
			errs = model.StepErrors{}
		}

		steps[i] = dto.FlowStepResponse{
			ID:          flow.Steps[i].StepID,
			Title:       flow.Steps[i].Title,
			Description: flow.Steps[i].Description,
			Status:      string(flow.Steps[i].Status),
			Skippable:   flow.Steps[i].Skippable,
			Data:        emptyObjectIfBlank(flow.Steps[i].Data),
			Errors:      errs,
		}
	}

	return &dto.FlowStateResponse{
		FlowType:    flow.FlowType,
		Status:      string(flow.Status),
		CurrentStep: flow.CurrentStepID,
		Steps:       steps,
		Metadata:    emptyObjectIfBlank(flow.Metadata),
	}
}
