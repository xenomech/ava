package routes

import (
	socketctrl "ava/api/internal/controller/socket"
	"ava/api/internal/middleware"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func socketRoutes(router fiber.Router, controller *socketctrl.Controller, mw *middleware.Middleware) {
	router.Get("/", mw.ValidateAccessToken, controller.Upgrade, websocket.New(controller.Serve))
}
