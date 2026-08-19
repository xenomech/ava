package routes

import (
	devicectrl "ava/api/internal/controller/device"
	"ava/api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func deviceRoutes(router fiber.Router, controller *devicectrl.Controller, mw *middleware.Middleware) {
	router.Get("/", mw.ValidateAccessToken, controller.List)
	router.Patch("/:deviceID", mw.ValidateAccessToken, controller.Rename)
	router.Post("/:deviceID/command", mw.ValidateAccessToken, controller.SendCommand)
}
