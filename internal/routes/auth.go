package routes

import (
	authctrl "ava/internal/controller/auth"
	"ava/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func authRoutes(router fiber.Router, controller *authctrl.Controller, authMiddleware *middleware.Middleware) {
	router.Post("/register", controller.Register)
	router.Post("/login", controller.Login)
}
