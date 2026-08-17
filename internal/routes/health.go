package routes

import (
	healthctrl "ava/internal/controller/health"

	"github.com/gofiber/fiber/v2"
)

func healthRoutes(router fiber.Router, controller *healthctrl.Controller) {
	router.Get("/", controller.Ping)
}
