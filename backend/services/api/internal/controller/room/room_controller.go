package room

import (
	"ava/api/internal/dto"
	roomsvc "ava/api/internal/services/room"
	"ava/api/pkg/response"
	"ava/api/pkg/serrors"
	"ava/api/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct {
	roomService roomsvc.Service
}

func NewController(roomService roomsvc.Service) *Controller {
	return &Controller{roomService: roomService}
}

func (c *Controller) List(ctx *fiber.Ctx) error {
	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	rooms, err := c.roomService.ListByTenant(ctx.Context(), tenantID)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to list rooms")
	}

	return response.Send(ctx, fiber.StatusOK, rooms, "")
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.CreateRoomRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	created, err := c.roomService.Create(ctx.Context(), tenantID, &req)
	if err != nil {
		return c.fail(ctx, err, "Failed to create the room")
	}

	return response.Send(ctx, fiber.StatusCreated, created, "")
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	roomID, err := uuid.Parse(ctx.Params("roomID"))
	if err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid room id")
	}

	var req dto.UpdateRoomRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	updated, err := c.roomService.Update(ctx.Context(), tenantID, roomID, &req)
	if err != nil {
		return c.fail(ctx, err, "Failed to update the room")
	}

	return response.Send(ctx, fiber.StatusOK, updated, "")
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	roomID, err := uuid.Parse(ctx.Params("roomID"))
	if err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid room id")
	}

	if err := c.roomService.Delete(ctx.Context(), tenantID, roomID); err != nil {
		return c.fail(ctx, err, "Failed to delete the room")
	}

	return response.Send(ctx, fiber.StatusNoContent, nil, "")
}

func (c *Controller) fail(ctx *fiber.Ctx, err error, fallback string) error {
	switch {
	case serrors.Is(err, roomsvc.ErrRoomNotFound):
		return response.SendError(ctx, fiber.StatusNotFound, err)
	case serrors.Is(err, roomsvc.ErrNameTaken):
		return response.SendError(ctx, fiber.StatusConflict, err)
	case serrors.Is(err, roomsvc.ErrNameRequired):
		return response.SendError(ctx, fiber.StatusUnprocessableEntity, err)
	default:
		return response.Send(ctx, fiber.StatusInternalServerError, nil, fallback)
	}
}
