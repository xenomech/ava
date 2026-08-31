package routes

import (
	roomctrl "ava/api/internal/controller/room"
	"ava/api/internal/middleware"
	"ava/api/internal/model"

	"github.com/gofiber/fiber/v2"
)

func roomRoutes(router fiber.Router, controller *roomctrl.Controller, mw *middleware.Middleware) {
	router.Get("/", mw.Authenticated, middleware.RequireScope(model.ScopeRoomsRead), controller.List)
	router.Post("/", mw.Authenticated, middleware.RequireScope(model.ScopeRoomsWrite), controller.Create)
	router.Patch(
		"/:roomID",
		mw.Authenticated,
		middleware.RequireScope(model.ScopeRoomsWrite),
		controller.Update,
	)
	router.Delete(
		"/:roomID",
		mw.Authenticated,
		middleware.RequireScope(model.ScopeRoomsWrite),
		controller.Delete,
	)
}
