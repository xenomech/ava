package flow

import (
	"context"
	"encoding/json"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	flowrepo "ava/api/internal/repository/flow"
	"ava/api/pkg/logger"
	"ava/api/pkg/serrors"

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

	return toFlowStateResponse(flow), nil
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

	metadataJSON := json.RawMessage("{}")

	if def.Metadata != nil {
		encoded, err := json.Marshal(def.Metadata)
		if err != nil {
			logger.Error("flow.createFlow.metadata", logger.Err(err))

			return nil, err
		}

		metadataJSON = encoded
	}

	flow := model.NewFlow(tenantID, userID, flowType, def.Steps[0].ID, metadataJSON)

	steps := make([]model.FlowStep, len(def.Steps))

	for i, stepDef := range def.Steps {
		status := model.FlowStepStatusPending
		if i == 0 {
			status = model.FlowStepStatusInProgress
		}

		steps[i] = model.NewFlowStep(
			tenantID,
			flow.ID,
			stepDef.ID,
			stepDef.Title,
			stepDef.Description,
			i,
			status,
			stepDef.Skippable,
		)
	}

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
