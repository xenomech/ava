package health

import (
	"ava/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) Ping(ctx *fiber.Ctx) error {
	data := c.healthSvc.HealthCheck()
	return response.Send(ctx, fiber.StatusOK, data, "")
}
