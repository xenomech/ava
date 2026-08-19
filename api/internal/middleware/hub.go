package middleware

import (
	"ava/api/internal/services/auth/jwt"
	hubsvc "ava/api/internal/services/hub"
	"ava/api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ValidateHubToken(hubService hubsvc.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractBearer(c)
		if token == "" {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Missing authorization token")
		}

		tokenManager := jwt.NewTokenManager()
		ctx := c.Context()

		claims, err := tokenManager.ValidateToken(token)
		if err != nil {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Invalid or expired token")
		}

		if claims.TokenType != jwt.HubToken {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Invalid token type")
		}

		if claims.HubID == uuid.Nil || claims.TenantID == uuid.Nil {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Token is not scoped to a hub")
		}

		hub, err := hubService.ValidateDevice(ctx, claims.HubID)
		if err != nil {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Hub has been revoked or removed")
		}

		if hub.TenantID != claims.TenantID {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Hub does not belong to this tenant")
		}

		c.Locals("hubID", hub.ID)
		c.Locals("tenantID", hub.TenantID)
		c.Locals("hub", hub)

		return c.Next()
	}
}
