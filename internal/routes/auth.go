package routes

import (
	authctrl "ava/internal/controller/auth"
	"ava/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func authRoutes(router fiber.Router, controller *authctrl.Controller, authMiddleware *middleware.Middleware) {
	router.Post("/register", controller.Register)
	router.Post("/verify-email", controller.VerifyEmail)
	router.Post("/resend-verification", controller.ResendVerification)
	router.Post("/login", controller.Login)
	router.Post("/refresh", controller.RefreshToken)
	router.Post("/forgot-password", controller.ForgotPassword)
	router.Post("/reset-password", controller.ResetPassword)
	router.Post("/accept-invite", controller.AcceptInvite)

	router.Get("/me", authMiddleware.ValidateAccessToken, controller.Me)
	router.Post("/switch-tenant", authMiddleware.ValidateAccessToken, controller.SwitchTenant)
	router.Post("/logout", authMiddleware.ValidateAccessToken, controller.Logout)
	router.Post("/logout-all", authMiddleware.ValidateAccessToken, controller.LogoutAll)
	router.Post("/change-password", authMiddleware.ValidateAccessToken, controller.ChangePassword)
}
