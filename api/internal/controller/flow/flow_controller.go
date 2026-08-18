package flow

import (
	"ava/api/internal/dto"
	flowsvc "ava/api/internal/services/flow"
	"ava/api/pkg/response"
	"ava/api/pkg/serrors"
	"ava/api/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) Get(ctx *fiber.Ctx) error {
	tenantID, userID, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	flow, err := c.flowService.GetFlow(ctx.Context(), tenantID, userID, ctx.Params("type"))
	if err != nil {
		if serrors.Is(err, flowsvc.ErrInvalidFlowType) {
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to get flow")
	}

	return response.Send(ctx, fiber.StatusOK, flow, "")
}

func (c *Controller) SubmitStep(ctx *fiber.Ctx) error {
	tenantID, userID, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.SubmitStepRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	flow, err := c.flowService.SubmitStep(ctx.Context(), tenantID, userID, ctx.Params("type"), ctx.Params("stepId"), &req)
	if err != nil {
		switch {
		case serrors.Is(err, flowsvc.ErrStepValidationFailed):
			return response.SendErrorWithData(ctx, fiber.StatusUnprocessableEntity, flow, err)
		case serrors.Is(err, flowsvc.ErrStepNotPermitted):
			return response.SendError(ctx, fiber.StatusForbidden, err)
		case serrors.Is(err, flowsvc.ErrStepNotCurrent), serrors.Is(err, flowsvc.ErrFlowAlreadyCompleted):
			return response.SendError(ctx, fiber.StatusConflict, err)
		case serrors.Is(err, flowsvc.ErrInvalidFlowType),
			serrors.Is(err, flowsvc.ErrStepNotFound),
			serrors.Is(err, flowsvc.ErrInvalidStepData):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to submit step")
	}

	return response.Send(ctx, fiber.StatusOK, flow, "")
}

func (c *Controller) GoBack(ctx *fiber.Ctx) error {
	tenantID, userID, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	flow, err := c.flowService.GoBack(ctx.Context(), tenantID, userID, ctx.Params("type"))
	if err != nil {
		switch {
		case serrors.Is(err, flowsvc.ErrNoPreviousStep), serrors.Is(err, flowsvc.ErrInvalidFlowType):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		case serrors.Is(err, flowsvc.ErrFlowAlreadyCompleted):
			return response.SendError(ctx, fiber.StatusConflict, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to go back")
	}

	return response.Send(ctx, fiber.StatusOK, flow, "")
}

func (c *Controller) SkipStep(ctx *fiber.Ctx) error {
	tenantID, userID, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	flow, err := c.flowService.SkipStep(ctx.Context(), tenantID, userID, ctx.Params("type"))
	if err != nil {
		switch {
		case serrors.Is(err, flowsvc.ErrStepNotSkippable), serrors.Is(err, flowsvc.ErrInvalidFlowType):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		case serrors.Is(err, flowsvc.ErrFlowAlreadyCompleted):
			return response.SendError(ctx, fiber.StatusConflict, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to skip step")
	}

	return response.Send(ctx, fiber.StatusOK, flow, "")
}
