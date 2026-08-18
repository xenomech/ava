package middleware

import (
	"time"

	"ava/api/internal/services"
	"ava/api/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Middleware struct {
	ValidateAccessToken fiber.Handler
	RequestTrace        fiber.Handler
}

func NewMiddleware(service *services.Service) *Middleware {
	return &Middleware{
		ValidateAccessToken: ValidateAccessToken(service.Auth),
		RequestTrace:        requestTrace(),
	}
}

func requestTrace() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		requestID, err := uuid.NewV7()
		if err != nil {
			requestID = uuid.New()
		}

		rid := requestID.String()

		c.Set("X-Request-ID", rid)
		c.Locals("requestID", rid)

		err = c.Next()

		logger.Info("REQUEST",
			logger.String("method", c.Method()),
			logger.String("path", c.Path()),
			logger.Int("status", c.Response().StatusCode()),
			logger.Duration("duration", time.Since(start)),
			logger.String("request_id", rid),
		)

		return err
	}
}
