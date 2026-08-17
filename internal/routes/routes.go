package routes

import (
	"ava/internal/controller"
	"ava/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func AddRoutes(app *fiber.App, ctrl *controller.Controller, mw *middleware.Middleware) {
	app.Use(middleware.GlobalRateLimit())

	api := app.Group("/api/v1")

	healthRoutes(api.Group("/health"), ctrl.Health)

	authRoutes(api.Group("/auth"), ctrl.Auth, mw)
	tenantRoutes(api.Group("/tenants"), ctrl.Tenant, mw)
	flowRoutes(api.Group("/flows"), ctrl.Flow, mw)
}
