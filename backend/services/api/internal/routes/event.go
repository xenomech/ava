package routes

import (
	eventctrl "ava/api/internal/controller/event"
	"ava/api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func eventRoutes(router fiber.Router, controller *eventctrl.Controller, mw *middleware.Middleware) {
	router.Get("/", mw.ValidateAccessToken, controller.Stream)
}
