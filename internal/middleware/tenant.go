package middleware

import (
	"slices"

	"ava/internal/model"
	"ava/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(roles ...model.TenantRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(model.TenantRole)
		if !ok {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Unauthorized")
		}

		if !slices.Contains(roles, role) {
			return response.Send(c, fiber.StatusForbidden, nil, "Insufficient permissions for this tenant")
		}

		return c.Next()
	}
}
