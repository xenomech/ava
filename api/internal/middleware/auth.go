package middleware

import (
	"strings"
	"time"

	authsvc "ava/api/internal/services/auth"
	"ava/api/internal/services/auth/jwt"
	"ava/api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ValidateAccessToken(authService authsvc.Service) fiber.Handler {
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

		if claims.TokenType != jwt.AccessToken {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Invalid token type")
		}

		if claims.TenantID == uuid.Nil {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Token is not scoped to a tenant")
		}

		session, err := authService.ValidateSession(ctx, claims.TenantID, claims.SessionID)
		if err != nil {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Session has been revoked or expired")
		}

		if time.Now().After(session.ExpiresAt) {
			return response.Send(c, fiber.StatusUnauthorized, nil, "Session has expired")
		}

		user, err := authService.GetUserByID(ctx, claims.UserID)
		if err != nil {
			return response.Send(c, fiber.StatusUnauthorized, nil, "User not found")
		}

		c.Locals("userID", claims.UserID)
		c.Locals("tenantID", claims.TenantID)
		c.Locals("role", claims.Role)
		c.Locals("sessionID", claims.SessionID)
		c.Locals("user", user)
		c.Locals("session", session)

		return c.Next()
	}
}

func extractToken(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}

	return parts[1]
}
