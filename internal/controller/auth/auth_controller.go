package auth

import (
	"ava/internal/dto"
	authsvc "ava/internal/services/auth"
	"ava/pkg/response"
	"ava/pkg/serrors"
	"ava/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *Controller) Register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	registered, err := c.authService.Register(ctx.Context(), &req)
	if err != nil {
		switch {
		case serrors.Is(err, authsvc.ErrInvalidSlug):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		case serrors.Is(err, authsvc.ErrUserAlreadyExists), serrors.Is(err, authsvc.ErrTenantAlreadyExists):
			return response.SendError(ctx, fiber.StatusConflict, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to register user")
	}

	return response.Send(ctx, fiber.StatusCreated, registered, "")
}

func (c *Controller) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	deviceInfo := dto.DeviceInfo{
		DeviceName: ctx.Get("X-Device-Name", "Unknown"),
		IPAddress:  ctx.IP(),
		UserAgent:  ctx.Get("User-Agent"),
	}

	authResponse, err := c.authService.Login(ctx.Context(), &req, deviceInfo)
	if err != nil {
		switch {
		case serrors.Is(err, authsvc.ErrInvalidCredentials):
			return response.SendError(ctx, fiber.StatusUnauthorized, err)
		case serrors.Is(err, authsvc.ErrNoTenantMembership), serrors.Is(err, authsvc.ErrAccessDenied):
			return response.SendError(ctx, fiber.StatusForbidden, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to login")
	}

	return response.Send(ctx, fiber.StatusOK, authResponse, "")
}

func (c *Controller) RefreshToken(ctx *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	tokens, err := c.authService.RefreshToken(ctx.Context(), req.RefreshToken)
	if err != nil {
		if serrors.Is(err, authsvc.ErrInvalidToken) ||
			serrors.Is(err, authsvc.ErrSessionRevoked) ||
			serrors.Is(err, authsvc.ErrSessionNotFound) ||
			serrors.Is(err, authsvc.ErrSessionExpired) {
			return response.SendError(ctx, fiber.StatusUnauthorized, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to refresh token")
	}

	return response.Send(ctx, fiber.StatusOK, tokens, "")
}

func (c *Controller) Logout(ctx *fiber.Ctx) error {
	sessionID, ok := ctx.Locals("sessionID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	if err := c.authService.Logout(ctx.Context(), tenantID, sessionID); err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to logout")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

func (c *Controller) Me(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	session, err := c.authService.CurrentSession(ctx.Context(), tenantID, userID)
	if err != nil {
		switch {
		case serrors.Is(err, authsvc.ErrAccessDenied):
			return response.SendError(ctx, fiber.StatusForbidden, err)
		case serrors.Is(err, authsvc.ErrUserNotFound):
			return response.SendError(ctx, fiber.StatusNotFound, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to load session")
	}

	return response.Send(ctx, fiber.StatusOK, session, "")
}

func (c *Controller) LogoutAll(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	if err := c.authService.LogoutAll(ctx.Context(), tenantID, userID); err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to logout")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

func (c *Controller) SwitchTenant(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	sessionID, ok := ctx.Locals("sessionID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.SwitchTenantRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	deviceInfo := dto.DeviceInfo{
		DeviceName: ctx.Get("X-Device-Name", "Unknown"),
		IPAddress:  ctx.IP(),
		UserAgent:  ctx.Get("User-Agent"),
	}

	authResponse, err := c.authService.SwitchTenant(ctx.Context(), tenantID, userID, sessionID, req.TenantSlug, deviceInfo)
	if err != nil {
		switch {
		case serrors.Is(err, authsvc.ErrAccessDenied):
			return response.SendError(ctx, fiber.StatusForbidden, err)
		case serrors.Is(err, authsvc.ErrUserNotFound):
			return response.SendError(ctx, fiber.StatusNotFound, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to switch tenant")
	}

	return response.Send(ctx, fiber.StatusOK, authResponse, "")
}
