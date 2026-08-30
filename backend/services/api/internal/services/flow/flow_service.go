package flow

import (
	"context"
	"encoding/json"
	"slices"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	flowrepo "ava/api/internal/repository/flow"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"

	"github.com/google/uuid"
)

func (s *flowService) GetFlow(ctx context.Context, tenantID, userID uuid.UUID, flowType string) (*dto.FlowStateResponse, error) {
	flow, err := s.flowRepo.GetByUserAndType(ctx, tenantID, userID, flowType)
	if err != nil {
		if serrors.Is(err, flowrepo.ErrFlowNotFound) {
			return s.createFlow(ctx, tenantID, userID, flowType)
		}

		logger.Error("flow.GetFlow", logger.Err(err))

		return nil, err
	}

	// Stored steps go stale when a definition changes, so re-cut them; a finished flow is left as the record.
	if flow.Status != model.FlowStatusCompleted && !matchesDefinition(flow, flowType) {
		return s.rebuildFlow(ctx, tenantID, flow, flowType)
	}

	return toFlowStateResponse(flow), nil
}

// rebuildFlow re-cuts an unfinished flow against the current definition, keeping the flow row itself.
func (s *flowService) rebuildFlow(ctx context.Context, tenantID uuid.UUID, flow *model.Flow, flowType string) (*dto.FlowStateResponse, error) {
	def, err := GetDefinition(flowType)
	if err != nil || len(def.Steps) == 0 {
		return nil, ErrInvalidFlowType
	}

	metadata, err := encodeMetadata(def)
	if err != nil {
		return nil, err
	}

	logger.Info("flow.rebuild",
		logger.String("flow.type", flowType),
		logger.Any("flow.ID", flow.ID),
	)

	flow.Status = model.FlowStatusInProgress
	flow.CurrentStepID = def.Steps[0].ID
	flow.Metadata = metadata

	if err := s.flowRepo.ReplaceSteps(ctx, tenantID, flow, buildSteps(tenantID, flow.ID, def)); err != nil {
		logger.Error("flow.rebuild.ReplaceSteps", logger.Err(err))

		return nil, err
	}

	return toFlowStateResponse(flow), nil
}

// matchesDefinition reports whether a stored flow still has exactly the registry's steps, in the same order.
func matchesDefinition(flow *model.Flow, flowType string) bool {
	def, err := GetDefinition(flowType)
	if err != nil {
		return false
	}

	if len(flow.Steps) != len(def.Steps) {
		return false
	}

	ordered := make([]model.FlowStep, len(flow.Steps))
	copy(ordered, flow.Steps)
	slices.SortFunc(ordered, func(a, b model.FlowStep) int { return a.StepOrder - b.StepOrder })

	for i := range def.Steps {
		if ordered[i].StepID != def.Steps[i].ID {
			return false
		}
	}

	return true
}

func (s *flowService) SubmitStep(ctx context.Context, tenantID, userID uuid.UUID, flowType, stepID string, req *dto.SubmitStepRequest) (*dto.FlowStateResponse, error) {
	flow, err := s.flowRepo.GetByUserAndType(ctx, tenantID, userID, flowType)
	if err != nil {
		if serrors.Is(err, flowrepo.ErrFlowNotFound) {
			return nil, ErrInvalidFlowType
		}

		logger.Error("flow.SubmitStep", logger.Err(err))

		return nil, err
	}

	if flow.Status == model.FlowStatusCompleted {
		return nil, ErrFlowAlreadyCompleted
	}

	if flow.CurrentStepID != stepID {
		return nil, ErrStepNotCurrent
	}

	currentStep := findStep(flow.Steps, stepID)
	if currentStep == nil {
		return nil, ErrStepNotFound
	}

	fieldErrors, err := s.runStepHandler(ctx, tenantID, userID, flowType, stepID, req.Data)
	if err != nil {
		return nil, err
	}

	if len(fieldErrors) > 0 {
		currentStep.Errors = fieldErrors
		currentStep.Status = model.FlowStepStatusFailed

		if err := s.flowRepo.UpdateStep(ctx, tenantID, currentStep); err != nil {
			logger.Error("flow.SubmitStep.SaveErrors", logger.Err(err))

			return nil, err
		}

		return toFlowStateResponse(flow), ErrStepValidationFailed
	}

	currentStep.Data = req.Data
	currentStep.Errors = model.StepErrors{}
	currentStep.Status = model.FlowStepStatusCompleted

	if err := s.flowRepo.UpdateStep(ctx, tenantID, currentStep); err != nil {
		logger.Error("flow.SubmitStep.UpdateStep", logger.Err(err))

		return nil, err
	}

	return s.advanceToNextStep(ctx, tenantID, flow)
}

func (s *flowService) runStepHandler(ctx context.Context, tenantID, userID uuid.UUID, flowType, stepID string, data json.RawMessage) (model.StepErrors, error) {
	handler, ok := s.handlers[flowType+":"+stepID]
	if !ok {
		return model.StepErrors{}, nil
	}

	fieldErrors, err := handler.Validate(ctx, tenantID, userID, data)
	if err != nil {
		return nil, logHandlerError("flow.SubmitStep.Validate", err)
	}

	if len(fieldErrors) > 0 {
		return fieldErrors, nil
	}

	if err := handler.Execute(ctx, tenantID, userID, data); err != nil {
		return nil, logHandlerError("flow.SubmitStep.Execute", err)
	}

	return model.StepErrors{}, nil
}

func logHandlerError(event string, err error) error {
	if serrors.Is(err, ErrStepNotPermitted) || serrors.Is(err, ErrInvalidStepData) {
		return err
	}

	logger.Error(event, logger.Err(err))

	return err
}

func (s *flowService) GoBack(ctx context.Context, tenantID, userID uuid.UUID, flowType string) (*dto.FlowStateResponse, error) {
	flow, err := s.flowRepo.GetByUserAndType(ctx, tenantID, userID, flowType)
	if err != nil {
		if serrors.Is(err, flowrepo.ErrFlowNotFound) {
			return nil, ErrInvalidFlowType
		}

		logger.Error("flow.GoBack", logger.Err(err))

		return nil, err
	}

	if flow.Status == model.FlowStatusCompleted {
		return nil, ErrFlowAlreadyCompleted
	}

	currentStep := findStep(flow.Steps, flow.CurrentStepID)
	if currentStep == nil {
		return nil, ErrStepNotFound
	}

	prevStep := findPreviousStep(flow.Steps, currentStep.StepOrder)
	if prevStep == nil {
		return nil, ErrNoPreviousStep
	}

	currentStep.Status = model.FlowStepStatusPending
	if err := s.flowRepo.UpdateStep(ctx, tenantID, currentStep); err != nil {
		logger.Error("flow.GoBack.ResetCurrent", logger.Err(err))

		return nil, err
	}

	prevStep.Status = model.FlowStepStatusInProgress
	if err := s.flowRepo.UpdateStep(ctx, tenantID, prevStep); err != nil {
		logger.Error("flow.GoBack.ActivatePrev", logger.Err(err))

		return nil, err
	}

	flow.CurrentStepID = prevStep.StepID
	flow.Status = model.FlowStatusInProgress

	if err := s.flowRepo.UpdateFlow(ctx, tenantID, flow); err != nil {
		logger.Error("flow.GoBack.UpdateFlow", logger.Err(err))

		return nil, err
	}

	return toFlowStateResponse(flow), nil
}

func (s *flowService) SkipStep(ctx context.Context, tenantID, userID uuid.UUID, flowType string) (*dto.FlowStateResponse, error) {
	flow, err := s.flowRepo.GetByUserAndType(ctx, tenantID, userID, flowType)
	if err != nil {
		if serrors.Is(err, flowrepo.ErrFlowNotFound) {
			return nil, ErrInvalidFlowType
		}

		logger.Error("flow.SkipStep", logger.Err(err))

		return nil, err
	}

	if flow.Status == model.FlowStatusCompleted {
		return nil, ErrFlowAlreadyCompleted
	}

	currentStep := findStep(flow.Steps, flow.CurrentStepID)
	if currentStep == nil {
		return nil, ErrStepNotFound
	}

	if !currentStep.Skippable {
		return nil, ErrStepNotSkippable
	}

	currentStep.Status = model.FlowStepStatusSkipped
	if err := s.flowRepo.UpdateStep(ctx, tenantID, currentStep); err != nil {
		logger.Error("flow.SkipStep.UpdateStep", logger.Err(err))

		return nil, err
	}

	return s.advanceToNextStep(ctx, tenantID, flow)
}

func (s *flowService) createFlow(ctx context.Context, tenantID, userID uuid.UUID, flowType string) (*dto.FlowStateResponse, error) {
	def, err := GetDefinition(flowType)
	if err != nil {
		return nil, ErrInvalidFlowType
	}

	if len(def.Steps) == 0 {
		return nil, ErrInvalidFlowType
	}

	metadataJSON, err := encodeMetadata(def)
	if err != nil {
		return nil, err
	}

	flow := model.NewFlow(tenantID, userID, flowType, def.Steps[0].ID, metadataJSON)
	steps := buildSteps(tenantID, flow.ID, def)

	if err := s.flowRepo.CreateFlowWithSteps(ctx, flow, steps); err != nil {
		if serrors.Is(err, flowrepo.ErrFlowAlreadyExists) {
			existing, err := s.flowRepo.GetByUserAndType(ctx, tenantID, userID, flowType)
			if err != nil {
				logger.Error("flow.createFlow.reread", logger.Err(err))

				return nil, err
			}

			return toFlowStateResponse(existing), nil
		}

		logger.Error("flow.createFlow", logger.Err(err))

		return nil, err
	}

	flow.Steps = steps

	return toFlowStateResponse(flow), nil
}

func (s *flowService) advanceToNextStep(ctx context.Context, tenantID uuid.UUID, flow *model.Flow) (*dto.FlowStateResponse, error) {
	nextStep := findNextPendingStep(flow.Steps)
	if nextStep == nil {
		flow.Status = model.FlowStatusCompleted
		flow.CurrentStepID = ""

		if err := s.flowRepo.UpdateFlow(ctx, tenantID, flow); err != nil {
			logger.Error("flow.advanceToNextStep.Complete", logger.Err(err))

			return nil, err
		}

		return toFlowStateResponse(flow), nil
	}

	nextStep.Status = model.FlowStepStatusInProgress
	if err := s.flowRepo.UpdateStep(ctx, tenantID, nextStep); err != nil {
		logger.Error("flow.advanceToNextStep.ActivateNext", logger.Err(err))

		return nil, err
	}

	flow.CurrentStepID = nextStep.StepID
	if err := s.flowRepo.UpdateFlow(ctx, tenantID, flow); err != nil {
		logger.Error("flow.advanceToNextStep.UpdateFlow", logger.Err(err))

		return nil, err
	}

	return toFlowStateResponse(flow), nil
}
