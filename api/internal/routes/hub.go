package routes

import (
	devicectrl "ava/api/internal/controller/device"
	hubctrl "ava/api/internal/controller/hub"
	"ava/api/internal/middleware"
	"ava/api/internal/model"

	"github.com/gofiber/fiber/v2"
)

func hubRoutes(router fiber.Router, controller *hubctrl.Controller, deviceController *devicectrl.Controller, mw *middleware.Middleware) {
	authRL := middleware.AuthRateLimit()

	router.Post("/code", authRL, controller.RequestCode)
	router.Post("/token", controller.Poll)
	router.Post("/token/refresh", controller.Refresh)

	router.Post("/heartbeat", mw.ValidateHubToken, controller.Heartbeat)
	router.Put("/devices", mw.ValidateHubToken, deviceController.Sync)

	router.Post("/activate", mw.ValidateAccessToken, controller.Activate)
	router.Get("/", mw.ValidateAccessToken, controller.List)
	router.Get("/:hubID/devices", mw.ValidateAccessToken, deviceController.ListByHub)
	router.Delete("/:hubID", mw.ValidateAccessToken,
		middleware.RequireRole(model.TenantRoleOwner, model.TenantRoleAdmin), controller.Revoke)
}
