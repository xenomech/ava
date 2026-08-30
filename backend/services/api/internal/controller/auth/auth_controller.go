package auth

import (
	"ava/api/internal/dto"
	authsvc "ava/api/internal/services/auth"
	"ava/api/pkg/response"
	"ava/api/pkg/serrors"
	"ava/api/pkg/validator"

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
		case serrors.Is(err, authsvc.ErrEmailNotVerified),
			serrors.Is(err, authsvc.ErrNoTenantMembership),
			serrors.Is(err, authsvc.ErrAccessDenied):
			return response.SendError(ctx, fiber.StatusForbidden, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to login")
	}

	setSessionCookies(ctx, authResponse.Tokens)

	authResponse.Tokens = redactTokens(authResponse.Tokens)

	return response.Send(ctx, fiber.StatusOK, authResponse, "")
}

func (c *Controller) RefreshToken(ctx *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	_ = ctx.BodyParser(&req)

	refreshToken := ctx.Cookies(RefreshCookie)
	if refreshToken == "" {
		refreshToken = req.RefreshToken
	}

	if refreshToken == "" {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Missing refresh token")
	}

	tokens, err := c.authService.RefreshToken(ctx.Context(), refreshToken)
	if err != nil {
		if serrors.Is(err, authsvc.ErrInvalidToken) ||
			serrors.Is(err, authsvc.ErrSessionRevoked) ||
			serrors.Is(err, authsvc.ErrSessionNotFound) ||
			serrors.Is(err, authsvc.ErrSessionExpired) {
			clearSessionCookies(ctx)

			return response.SendError(ctx, fiber.StatusUnauthorized, err)
		}

		clearSessionCookies(ctx)

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to refresh token")
	}

	setSessionCookies(ctx, tokens)

	return response.Send(ctx, fiber.StatusOK, redactTokens(tokens), "")
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

	clearSessionCookies(ctx)

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

	clearSessionCookies(ctx)

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

	setSessionCookies(ctx, authResponse.Tokens)

	authResponse.Tokens = redactTokens(authResponse.Tokens)

	return response.Send(ctx, fiber.StatusOK, authResponse, "")
}

func (c *Controller) AcceptInvite(ctx *fiber.Ctx) error {
	var req dto.AcceptInviteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	tenant, err := c.authService.AcceptInvite(ctx.Context(), req.Token)
	if err != nil {
		if serrors.Is(err, authsvc.ErrInviteInvalid) {
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to accept invitation")
	}

	return response.Send(ctx, fiber.StatusOK, tenant, "")
}

func (c *Controller) VerifyEmail(ctx *fiber.Ctx) error {
	var req dto.VerifyEmailRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	if err := c.authService.VerifyEmail(ctx.Context(), req.Token); err != nil {
		if serrors.Is(err, authsvc.ErrInvalidToken) {
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to verify email")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

func (c *Controller) ResendVerification(ctx *fiber.Ctx) error {
	var req dto.ResendVerificationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	if err := c.authService.ResendVerification(ctx.Context(), req.Email); err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to resend verification email")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

func (c *Controller) ForgotPassword(ctx *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	if err := c.authService.ForgotPassword(ctx.Context(), req.Email); err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to process request")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

func (c *Controller) ResetPassword(ctx *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	if err := c.authService.ResetPassword(ctx.Context(), req.Token, req.NewPassword); err != nil {
		if serrors.Is(err, authsvc.ErrInvalidToken) {
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to reset password")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

func (c *Controller) ChangePassword(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	var req dto.ChangePasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Send(ctx, fiber.StatusBadRequest, nil, "Invalid request body")
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return response.SendValidation(ctx, validator.FirstError(err), validator.FieldErrors(err))
	}

	if err := c.authService.ChangePassword(ctx.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		switch {
		case serrors.Is(err, authsvc.ErrPasswordMismatch):
			return response.SendError(ctx, fiber.StatusBadRequest, err)
		case serrors.Is(err, authsvc.ErrUserNotFound):
			return response.SendError(ctx, fiber.StatusNotFound, err)
		}

		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to change password")
	}

	return response.Send(ctx, fiber.StatusOK, nil, "")
}

func (c *Controller) GetSessions(ctx *fiber.Ctx) error {
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

	sessions, err := c.authService.GetUserSessions(ctx.Context(), tenantID, userID, sessionID)
	if err != nil {
		return response.Send(ctx, fiber.StatusInternalServerError, nil, "Failed to load sessions")
	}

	return response.Send(ctx, fiber.StatusOK, sessions, "")
}
