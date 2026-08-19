package middleware

import (
	"slices"
	"strings"

	"ava/api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func VerifyOrigin(allowedOrigins string) fiber.Handler {
	allowed := parseOrigins(allowedOrigins)

	return func(c *fiber.Ctx) error {
		switch c.Method() {
		case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions:
			return c.Next()
		}

		origin := c.Get("Origin")
		if origin == "" {
			return c.Next()
		}

		if slices.Contains(allowed, "*") || slices.Contains(allowed, origin) {
			return c.Next()
		}

		return response.Send(c, fiber.StatusForbidden, nil, "Origin not allowed")
	}
}

func parseOrigins(value string) []string {
	parts := strings.Split(value, ",")

	origins := make([]string, 0, len(parts))

	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	return origins
}
