package routes

import (
	scenectrl "ava/api/internal/controller/scene"
	"ava/api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func sceneRoutes(router fiber.Router, controller *scenectrl.Controller, mw *middleware.Middleware) {
	router.Get("/", mw.ValidateAccessToken, controller.List)
	router.Post("/", mw.ValidateAccessToken, controller.Create)
	router.Delete("/:sceneID", mw.ValidateAccessToken, controller.Delete)
}
