package routes

import (
	roomctrl "ava/api/internal/controller/room"
	"ava/api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func roomRoutes(router fiber.Router, controller *roomctrl.Controller, mw *middleware.Middleware) {
	router.Get("/", mw.ValidateAccessToken, controller.List)
	router.Post("/", mw.ValidateAccessToken, controller.Create)
	router.Patch("/:roomID", mw.ValidateAccessToken, controller.Update)
	router.Delete("/:roomID", mw.ValidateAccessToken, controller.Delete)
}
