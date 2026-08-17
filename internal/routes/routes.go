package routes

import (
	"ava/internal/controller"
	"ava/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func AddRoutes(app *fiber.App, ctrl *controller.Controller, mw *middleware.Middleware) {
	api := app.Group("/api/v1")

	healthRoutes(api.Group("/health"), ctrl.Health)

	authRoutes(api.Group("/auth"), ctrl.Auth, mw)
	tenantRoutes(api.Group("/tenants"), ctrl.Tenant, mw)
}
