package apitoken

import (
	"context"
	"time"

	"ava/api/internal/dto"
	"ava/api/internal/middleware"
	"ava/api/internal/model"
	apitokensvc "ava/api/internal/services/apitoken"
	"ava/api/pkg/response"
	"ava/api/pkg/serrors"
	"ava/api/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Scopes lists every scope this server understands, so a client need not hardcode them.
func (c *Controller) Scopes(ctx *fiber.Ctx) error {
	names := make([]string, 0, len(model.AllScopes))
	for _, scope := range model.AllScopes {
		names = append(names, string(scope))
	}

	return response.Send(ctx, fiber.StatusOK, dto.ScopeResponse{Scopes: names}, "")
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	var req dto.CreateAPITokenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	tenantID, userID, ok := middleware.Actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Not authenticated")
	}

	var expiresAt *time.Time
	if req.ExpiresInDays > 0 {
		at := time.Now().AddDate(0, 0, req.ExpiresInDays)
		expiresAt = &at
	}

	token, plaintext, err := c.apiTokenService.Create(
		ctx.Context(), tenantID, userID, req.Name, req.Scopes, expiresAt,
	)
	if err != nil {
		switch {
		case serrors.Is(err, apitokensvc.ErrNameTaken):
			return response.SendError(ctx, fiber.StatusConflict, err)
		case serrors.Is(err, apitokensvc.ErrUnknownScope), serrors.Is(err, apitokensvc.ErrNoScopes):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to create token")
	}

	return response.Send(ctx, fiber.StatusCreated, dto.CreatedAPITokenResponse{
		Token: dto.NewAPITokenResponse(token),
		Value: plaintext,
	}, "")
}

func (c *Controller) List(ctx *fiber.Ctx) error {
	tenantID, userID, ok := middleware.Actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Not authenticated")
	}

	tokens, err := c.apiTokenService.List(ctx.Context(), tenantID, userID)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to list tokens")
	}

	out := make([]dto.APITokenResponse, 0, len(tokens))
	for _, token := range tokens {
		out = append(out, dto.NewAPITokenResponse(token))
	}

	return response.Send(ctx, fiber.StatusOK, out, "")
}

func (c *Controller) Revoke(ctx *fiber.Ctx) error {
	return c.mutate(ctx, c.apiTokenService.Revoke, "Failed to revoke token")
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	return c.mutate(ctx, c.apiTokenService.Delete, "Failed to delete token")
}

func (c *Controller) mutate(
	ctx *fiber.Ctx,
	apply func(reqCtx context.Context, tenantID, userID, tokenID uuid.UUID) error,
	failure string,
) error {
	tenantID, userID, ok := middleware.Actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Not authenticated")
	}

	tokenID, err := uuid.Parse(ctx.Params("tokenID"))
	if err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid token id")
	}

	if err := apply(ctx.Context(), tenantID, userID, tokenID); err != nil {
		if serrors.Is(err, apitokensvc.ErrTokenNotFound) {
			return response.SendError(ctx, fiber.StatusNotFound, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, failure)
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}
