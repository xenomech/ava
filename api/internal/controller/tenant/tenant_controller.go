package tenant

import (
	"ava/api/internal/dto"
	"ava/api/internal/response"
	"ava/api/internal/serrors"
	tenantsvc "ava/api/internal/services/tenant"
	"ava/api/internal/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *Controller) Create(ctx *fiber.Ctx) error {
	_, userID, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.CreateTenantRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	created, err := c.tenantService.Create(ctx.Context(), userID, &req)
	if err != nil {
		switch {
		case serrors.Is(err, tenantsvc.ErrInvalidSlug):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		case serrors.Is(err, tenantsvc.ErrTenantAlreadyExists):
			return response.SendError(ctx, fiber.StatusConflict, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to create tenant")
	}

	return response.Send(ctx, fiber.StatusCreated, created, "")
}

func (c *Controller) ListMine(ctx *fiber.Ctx) error {
	_, userID, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	tenants, err := c.tenantService.ListMine(ctx.Context(), userID)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to list tenants")
	}

	return response.Send(ctx, fiber.StatusOK, tenants, "")
}

func (c *Controller) Get(ctx *fiber.Ctx) error {
	tenantID, _, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	found, err := c.tenantService.Get(ctx.Context(), tenantID)
	if err != nil {
		if serrors.Is(err, tenantsvc.ErrTenantNotFound) {
			return response.SendError(ctx, fiber.StatusNotFound, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to load tenant")
	}

	return response.Send(ctx, fiber.StatusOK, found, "")
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	tenantID, _, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.UpdateTenantRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	updated, err := c.tenantService.Update(ctx.Context(), tenantID, &req)
	if err != nil {
		if serrors.Is(err, tenantsvc.ErrTenantNotFound) {
			return response.SendError(ctx, fiber.StatusNotFound, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to update tenant")
	}

	return response.Send(ctx, fiber.StatusOK, updated, "")
}

func (c *Controller) ListMembers(ctx *fiber.Ctx) error {
	tenantID, _, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	members, err := c.tenantService.ListMembers(ctx.Context(), tenantID)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to list members")
	}

	return response.Send(ctx, fiber.StatusOK, members, "")
}

func (c *Controller) UpdateMemberRole(ctx *fiber.Ctx) error {
	tenantID, _, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	memberID, err := uuid.Parse(ctx.Params("userID"))
	if err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid user ID")
	}

	var req dto.UpdateMemberRoleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	if err := c.tenantService.UpdateMemberRole(ctx.Context(), tenantID, memberID, &req); err != nil {
		switch {
		case serrors.Is(err, tenantsvc.ErrInvalidRole):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		case serrors.Is(err, tenantsvc.ErrMemberNotFound):
			return response.SendError(ctx, fiber.StatusNotFound, err)
		case serrors.Is(err, tenantsvc.ErrLastOwner):
			return response.SendError(ctx, fiber.StatusConflict, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to update member role")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

func (c *Controller) RemoveMember(ctx *fiber.Ctx) error {
	tenantID, _, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	memberID, err := uuid.Parse(ctx.Params("userID"))
	if err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid user ID")
	}

	if err := c.tenantService.RemoveMember(ctx.Context(), tenantID, memberID); err != nil {
		switch {
		case serrors.Is(err, tenantsvc.ErrMemberNotFound):
			return response.SendError(ctx, fiber.StatusNotFound, err)
		case serrors.Is(err, tenantsvc.ErrLastOwner):
			return response.SendError(ctx, fiber.StatusConflict, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to remove member")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

func (c *Controller) Invite(ctx *fiber.Ctx) error {
	tenantID, userID, ok := actor(ctx)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.InviteMemberRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	invited, err := c.tenantService.Invite(ctx.Context(), tenantID, userID, &req)
	if err != nil {
		switch {
		case serrors.Is(err, tenantsvc.ErrInvalidRole):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		case serrors.Is(err, tenantsvc.ErrUserNotFound):
			return response.SendError(ctx, fiber.StatusNotFound, err)
		case serrors.Is(err, tenantsvc.ErrAlreadyMember), serrors.Is(err, tenantsvc.ErrAlreadyInvited):
			return response.SendError(ctx, fiber.StatusConflict, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to invite member")
	}

	return response.Send(ctx, fiber.StatusCreated, invited, "")
}
