package routes

import (
	tenantctrl "ava/internal/controller/tenant"
	"ava/internal/middleware"
	"ava/internal/model"

	"github.com/gofiber/fiber/v2"
)

func tenantRoutes(router fiber.Router, controller *tenantctrl.Controller, mw *middleware.Middleware) {
	router.Use(mw.ValidateAccessToken)

	router.Get("/", controller.ListMine)
	router.Post("/", controller.Create)

	router.Get("/current", controller.Get)
	router.Patch("/current", middleware.RequireRole(model.TenantRoleOwner, model.TenantRoleAdmin), controller.Update)

	router.Get("/current/members", controller.ListMembers)
	router.Patch("/current/members/:userID", middleware.RequireRole(model.TenantRoleOwner), controller.UpdateMemberRole)
	router.Delete("/current/members/:userID", middleware.RequireRole(model.TenantRoleOwner, model.TenantRoleAdmin), controller.RemoveMember)
}
