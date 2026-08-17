package auth

import (
	"ava/internal/dto"
	authsvc "ava/internal/services/auth"
	"ava/pkg/response"
	"ava/pkg/serrors"
	"ava/pkg/validator"

	"github.com/gofiber/fiber/v2"
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
		if serrors.Is(err, authsvc.ErrUserAlreadyExists) {
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
		if serrors.Is(err, authsvc.ErrInvalidCredentials) {
			return response.SendError(ctx, fiber.StatusUnauthorized, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to login")
	}

	return response.Send(ctx, fiber.StatusOK, authResponse, "")
}
