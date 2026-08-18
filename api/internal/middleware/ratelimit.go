package middleware

import (
	"time"

	"ava/api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

const CodeAuthRateLimited = "auth_rate_limited"

func GlobalRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               100,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return response.SendWithCode(c, fiber.StatusTooManyRequests, nil,
				response.CodeRateLimited, "Too many requests, please try again later")
		},
	})
}

func AuthRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               25,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return response.SendWithCode(c, fiber.StatusTooManyRequests, nil,
				CodeAuthRateLimited, "Too many authentication attempts, please try again later")
		},
	})
}
