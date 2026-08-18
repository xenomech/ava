package routes

import (
	devicectrl "ava/api/internal/controller/device"
	"ava/api/internal/middleware"
	"ava/api/internal/model"

	"github.com/gofiber/fiber/v2"
)

func deviceRoutes(router fiber.Router, controller *devicectrl.Controller, mw *middleware.Middleware) {
	authRL := middleware.AuthRateLimit()

	router.Post("/code", authRL, controller.RequestCode)
	router.Post("/token", controller.Poll)
	router.Post("/token/refresh", controller.Refresh)

	router.Post("/heartbeat", mw.ValidateDeviceToken, controller.Heartbeat)

	router.Post("/activate", mw.ValidateAccessToken, controller.Activate)
	router.Get("/", mw.ValidateAccessToken, controller.List)
	router.Delete("/:deviceID", mw.ValidateAccessToken,
		middleware.RequireRole(model.TenantRoleOwner, model.TenantRoleAdmin), controller.Revoke)
}
