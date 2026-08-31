package routes

import (
	devicectrl "ava/api/internal/controller/device"
	"ava/api/internal/middleware"
	"ava/api/internal/model"

	"github.com/gofiber/fiber/v2"
)

func deviceRoutes(router fiber.Router, controller *devicectrl.Controller, mw *middleware.Middleware) {
	router.Get("/", mw.Authenticated, middleware.RequireScope(model.ScopeDevicesRead), controller.List)
	router.Patch(
		"/:deviceID",
		mw.Authenticated,
		middleware.RequireScope(model.ScopeDevicesWrite),
		controller.Update,
	)
	router.Post(
		"/apply",
		mw.Authenticated,
		middleware.RequireScope(model.ScopeDevicesWrite),
		controller.Apply,
	)
	router.Post(
		"/:deviceID/command",
		mw.Authenticated,
		middleware.RequireScope(model.ScopeDevicesWrite),
		controller.SendCommand,
	)
}
