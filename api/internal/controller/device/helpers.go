package device

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func actor(ctx *fiber.Ctx) (tenantID, userID uuid.UUID, ok bool) {
	tenantID, ok = ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}

	userID, ok = ctx.Locals("userID").(uuid.UUID)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}

	return tenantID, userID, true
}
