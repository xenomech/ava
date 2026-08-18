package middleware

import (
	"ava/api/internal/services/auth/jwt"
	devicesvc "ava/api/internal/services/device"
	"ava/api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ValidateDeviceToken(deviceService devicesvc.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractToken(c)
		if token == "" {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Missing authorization token")
		}

		tokenManager := jwt.NewTokenManager()
		ctx := c.Context()

		claims, err := tokenManager.ValidateToken(token)
		if err != nil {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Invalid or expired token")
		}

		if claims.TokenType != jwt.DeviceToken {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Invalid token type")
		}

		if claims.DeviceID == uuid.Nil || claims.TenantID == uuid.Nil {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Token is not scoped to a device")
		}

		device, err := deviceService.ValidateDevice(ctx, claims.DeviceID)
		if err != nil {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Device has been revoked or removed")
		}

		if device.TenantID != claims.TenantID {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Device does not belong to this tenant")
		}

		c.Locals("deviceID", device.ID)
		c.Locals("tenantID", device.TenantID)
		c.Locals("device", device)

		return c.Next()
	}
}
