package middleware

import (
	"ava/internal/services"

	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	ValidateAccessToken fiber.Handler
}

func NewMiddleware(service *services.Service) *Middleware {
	return &Middleware{
		ValidateAccessToken: ValidateAccessToken(service.Auth),
	}
}
