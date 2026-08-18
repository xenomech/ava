package device

import (
	"ava/api/internal/dto"
	devicesvc "ava/api/internal/services/device"
	"ava/api/pkg/response"
	"ava/api/pkg/serrors"
	"ava/api/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *Controller) RequestCode(ctx *fiber.Ctx) error {
	var req dto.DeviceCodeRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	code, err := c.deviceService.RequestCode(ctx.Context(), &req)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to start device authorization")
	}

	return response.Send(ctx, fiber.StatusCreated, code, "")
}

func (c *Controller) Poll(ctx *fiber.Ctx) error {
	var req dto.DevicePollRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	tokens, err := c.deviceService.Poll(ctx.Context(), req.DeviceCode)
	if err != nil {
		switch {
		case serrors.Is(err, devicesvc.ErrAuthorizationPending),
			serrors.Is(err, devicesvc.ErrExpiredCode),
			serrors.Is(err, devicesvc.ErrInvalidCode):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		case serrors.Is(err, devicesvc.ErrSlowDown):
			return response.SendError(ctx, fiber.StatusTooManyRequests, err)
		case serrors.Is(err, devicesvc.ErrAccessDenied):
			return response.SendError(ctx, fiber.StatusForbidden, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to poll device authorization")
	}

	return response.Send(ctx, fiber.StatusOK, tokens, "")
}

func (c *Controller) Refresh(ctx *fiber.Ctx) error {
	var req dto.DeviceRefreshRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	tokens, err := c.deviceService.Refresh(ctx.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case serrors.Is(err, devicesvc.ErrInvalidRefreshToken):
			return response.SendError(ctx, fiber.StatusUnauthorized, err)
		case serrors.Is(err, devicesvc.ErrDeviceRevoked):
			return response.SendError(ctx, fiber.StatusForbidden, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to refresh device token")
	}

	return response.Send(ctx, fiber.StatusOK, tokens, "")
}

func (c *Controller) Activate(ctx *fiber.Ctx) error {
	tenantID, userID, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.ActivateDeviceRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	device, err := c.deviceService.Activate(ctx.Context(), tenantID, userID, req.UserCode)
	if err != nil {
		switch {
		case serrors.Is(err, devicesvc.ErrInvalidCode):
			return response.SendError(ctx, fiber.StatusNotFound, err)
		case serrors.Is(err, devicesvc.ErrExpiredCode):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to activate device")
	}

	return response.Send(ctx, fiber.StatusCreated, device, "")
}

func (c *Controller) List(ctx *fiber.Ctx) error {
	tenantID, _, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	devices, err := c.deviceService.ListByTenant(ctx.Context(), tenantID)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to list devices")
	}

	return response.Send(ctx, fiber.StatusOK, devices, "")
}

func (c *Controller) Revoke(ctx *fiber.Ctx) error {
	tenantID, _, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	deviceID, err := uuid.Parse(ctx.Params("deviceID"))
	if err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid device ID")
	}

	if err := c.deviceService.Revoke(ctx.Context(), tenantID, deviceID); err != nil {
		if serrors.Is(err, devicesvc.ErrDeviceNotFound) {
			return response.SendError(ctx, fiber.StatusNotFound, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to revoke device")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

func (c *Controller) Heartbeat(ctx *fiber.Ctx) error {
	deviceID, ok := ctx.Locals("deviceID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	if err := c.deviceService.Heartbeat(ctx.Context(), deviceID); err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to record heartbeat")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}
