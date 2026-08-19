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

func (c *Controller) Sync(ctx *fiber.Ctx) error {
	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	hubID, ok := ctx.Locals("hubID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.SyncDevicesRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	devices, err := c.deviceService.SyncFromHub(ctx.Context(), tenantID, hubID, &req)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to sync devices")
	}

	return response.Send(ctx, fiber.StatusOK, devices, "")
}

func (c *Controller) List(ctx *fiber.Ctx) error {
	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	devices, err := c.deviceService.ListByTenant(ctx.Context(), tenantID)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to list devices")
	}

	return response.Send(ctx, fiber.StatusOK, devices, "")
}

func (c *Controller) ListByHub(ctx *fiber.Ctx) error {
	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	hubID, err := uuid.Parse(ctx.Params("hubID"))
	if err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid hub ID")
	}

	devices, err := c.deviceService.ListByHub(ctx.Context(), tenantID, hubID)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to list devices")
	}

	return response.Send(ctx, fiber.StatusOK, devices, "")
}

func (c *Controller) Rename(ctx *fiber.Ctx) error {
	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	deviceID, err := uuid.Parse(ctx.Params("deviceID"))
	if err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid device ID")
	}

	var req dto.RenameDeviceRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	device, err := c.deviceService.Rename(ctx.Context(), tenantID, deviceID, req.Name)
	if err != nil {
		if serrors.Is(err, devicesvc.ErrDeviceNotFound) {
			return response.SendError(ctx, fiber.StatusNotFound, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to rename device")
	}

	return response.Send(ctx, fiber.StatusOK, device, "")
}
