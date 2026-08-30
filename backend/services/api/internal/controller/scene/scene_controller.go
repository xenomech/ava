package scene

import (
	"ava/api/internal/dto"
	scenesvc "ava/api/internal/services/scene"
	"ava/api/pkg/response"
	"ava/api/pkg/serrors"
	"ava/api/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct {
	sceneService scenesvc.Service
}

func NewController(sceneService scenesvc.Service) *Controller {
	return &Controller{sceneService: sceneService}
}

func (c *Controller) List(ctx *fiber.Ctx) error {
	tenantID, roomID, err := c.scope(ctx)
	if err != nil {
		return err
	}

	scenes, err := c.sceneService.ListByRoom(ctx.Context(), tenantID, roomID)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to list scenes")
	}

	return response.Send(ctx, fiber.StatusOK, scenes, "")
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	tenantID, roomID, err := c.scope(ctx)
	if err != nil {
		return err
	}

	var req dto.CreateSceneRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	created, err := c.sceneService.Create(ctx.Context(), tenantID, roomID, &req)
	if err != nil {
		return c.fail(ctx, err, "Failed to save the scene")
	}

	return response.Send(ctx, fiber.StatusCreated, created, "")
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	tenantID, roomID, err := c.scope(ctx)
	if err != nil {
		return err
	}

	sceneID, parseErr := uuid.Parse(ctx.Params("sceneID"))
	if parseErr != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid scene id")
	}

	if err := c.sceneService.Delete(ctx.Context(), tenantID, roomID, sceneID); err != nil {
		return c.fail(ctx, err, "Failed to delete the scene")
	}

	return response.Send(ctx, fiber.StatusNoContent, nil, "")
}

// scope reads the two identifiers every route here needs, and answers the
// request itself if either is missing.
func (c *Controller) scope(ctx *fiber.Ctx) (tenantID, roomID uuid.UUID, err error) {
	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return uuid.Nil, uuid.Nil, response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	roomID, parseErr := uuid.Parse(ctx.Params("roomID"))
	if parseErr != nil {
		return uuid.Nil, uuid.Nil, response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid room id")
	}

	return tenantID, roomID, nil
}

func (c *Controller) fail(ctx *fiber.Ctx, err error, fallback string) error {
	switch {
	case serrors.Is(err, scenesvc.ErrSceneNotFound), serrors.Is(err, scenesvc.ErrRoomNotFound):
		return response.SendError(ctx, fiber.StatusNotFound, err)
	case serrors.Is(err, scenesvc.ErrNameTaken):
		return response.SendError(ctx, fiber.StatusConflict, err)
	case serrors.Is(err, scenesvc.ErrNameRequired), serrors.Is(err, scenesvc.ErrNothingToSave):
		return response.SendError(ctx, fiber.StatusUnprocessableEntity, err)
	default:
		return response.Send(ctx, fiber.StatusInternalServerError, nil, fallback)
	}
}
