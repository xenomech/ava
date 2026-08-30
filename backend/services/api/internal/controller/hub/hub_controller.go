package hub

import (
	"ava/api/internal/dto"
	hubsvc "ava/api/internal/services/hub"
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

	code, err := c.hubService.RequestCode(ctx.Context(), &req)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to start hub authorization")
	}

	return response.Send(ctx, fiber.StatusCreated, code, "")
}

func (c *Controller) Poll(ctx *fiber.Ctx) error {
	var req dto.HubPollRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	tokens, err := c.hubService.Poll(ctx.Context(), req.DeviceCode)
	if err != nil {
		switch {
		case serrors.Is(err, hubsvc.ErrAuthorizationPending),
			serrors.Is(err, hubsvc.ErrExpiredCode),
			serrors.Is(err, hubsvc.ErrInvalidCode):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		case serrors.Is(err, hubsvc.ErrSlowDown):
			return response.SendError(ctx, fiber.StatusTooManyRequests, err)
		case serrors.Is(err, hubsvc.ErrAccessDenied):
			return response.SendError(ctx, fiber.StatusForbidden, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to poll hub authorization")
	}

	return response.Send(ctx, fiber.StatusOK, tokens, "")
}

func (c *Controller) Refresh(ctx *fiber.Ctx) error {
	var req dto.HubRefreshRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	tokens, err := c.hubService.Refresh(ctx.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case serrors.Is(err, hubsvc.ErrInvalidRefreshToken):
			return response.SendError(ctx, fiber.StatusUnauthorized, err)
		case serrors.Is(err, hubsvc.ErrHubRevoked):
			return response.SendError(ctx, fiber.StatusForbidden, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to refresh hub token")
	}

	return response.Send(ctx, fiber.StatusOK, tokens, "")
}

func (c *Controller) Activate(ctx *fiber.Ctx) error {
	tenantID, userID, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.ActivateHubRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	hub, err := c.hubService.Activate(ctx.Context(), tenantID, userID, req.UserCode)
	if err != nil {
		switch {
		case serrors.Is(err, hubsvc.ErrInvalidCode):
			return response.SendError(ctx, fiber.StatusNotFound, err)
		case serrors.Is(err, hubsvc.ErrExpiredCode):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to activate hub")
	}

	return response.Send(ctx, fiber.StatusCreated, hub, "")
}

func (c *Controller) List(ctx *fiber.Ctx) error {
	tenantID, _, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	hubs, err := c.hubService.ListByTenant(ctx.Context(), tenantID)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to list hubs")
	}

	return response.Send(ctx, fiber.StatusOK, hubs, "")
}

func (c *Controller) Revoke(ctx *fiber.Ctx) error {
	tenantID, _, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	hubID, err := uuid.Parse(ctx.Params("hubID"))
	if err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid hub ID")
	}

	if err := c.hubService.Revoke(ctx.Context(), tenantID, hubID); err != nil {
		if serrors.Is(err, hubsvc.ErrHubNotFound) {
			return response.SendError(ctx, fiber.StatusNotFound, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to revoke hub")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

// Heartbeat records that a hub is alive, and whether it can still be reached.
//
// Those are not the same thing, and treating them as one is what let a hub sit
// there looking healthy while every command vanished. Heartbeats and device
// syncs are HTTP; commands travel over the broker. When the broker stopped
// accepting the hub, the HTTP half carried on perfectly and nothing noticed —
// the API went on publishing into a topic with no subscriber and reporting
// success, because publishing had in fact succeeded.
//
// The hub is the only party that can know whether its command channel is up, so
// it is asked to say.
func (c *Controller) Heartbeat(ctx *fiber.Ctx) error {
	hubID, ok := ctx.Locals("hubID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.HubHeartbeatRequest
	if err := ctx.BodyParser(&req); err != nil {
		req = dto.HubHeartbeatRequest{}
	}

	if req.BrokerConnected != nil {
		if err := c.hubService.ApplyPresence(ctx.Context(), hubID, *req.BrokerConnected); err != nil {
			return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to record presence")
		}
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}
