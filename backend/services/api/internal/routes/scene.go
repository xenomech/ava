package routes

import (
	scenectrl "ava/api/internal/controller/scene"
	"ava/api/internal/middleware"
	"ava/api/internal/model"

	"github.com/gofiber/fiber/v2"
)

func sceneRoutes(router fiber.Router, controller *scenectrl.Controller, mw *middleware.Middleware) {
	router.Get("/", mw.Authenticated, middleware.RequireScope(model.ScopeScenesRead), controller.List)
	router.Post(
		"/",
		mw.Authenticated,
		middleware.RequireScope(model.ScopeScenesWrite),
		controller.Create,
	)
	// Applying reads the scene and writes the devices, so it needs a scope for each.
	router.Post(
		"/:sceneID/apply",
		mw.Authenticated,
		middleware.RequireScope(model.ScopeScenesRead),
		middleware.RequireScope(model.ScopeDevicesWrite),
		controller.Apply,
	)
	router.Delete(
		"/:sceneID",
		mw.Authenticated,
		middleware.RequireScope(model.ScopeScenesWrite),
		controller.Delete,
	)
}
