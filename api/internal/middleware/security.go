package middleware

import "github.com/gofiber/fiber/v2"

func SecurityHeaders(serverEnv string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "0")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		if serverEnv != "local" {
			c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		return c.Next()
	}
}
