package event

import (
	"bufio"
	"fmt"
	"time"

	eventsvc "ava/api/internal/services/event"
	"ava/api/pkg/response"
	"ava/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

const heartbeatInterval = 25 * time.Second

type Controller struct {
	events eventsvc.Service
}

func NewController(events eventsvc.Service) *Controller {
	return &Controller{events: events}
}

func (c *Controller) Stream(ctx *fiber.Ctx) error {
	tenantID, ok := ctx.Locals("tenantID").(uuid.UUID)
	if !ok {
		return response.Send(ctx, fiber.StatusUnauthorized, nil, "Unauthorized")
	}

	ctx.Set(fiber.HeaderContentType, "text/event-stream")
	ctx.Set(fiber.HeaderCacheControl, "no-cache")
	ctx.Set(fiber.HeaderConnection, "keep-alive")
	ctx.Set("X-Accel-Buffering", "no")

	stream, cancel := c.events.Subscribe(tenantID)

	ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(writer *bufio.Writer) {
		defer cancel()

		heartbeat := time.NewTicker(heartbeatInterval)
		defer heartbeat.Stop()

		if _, err := fmt.Fprint(writer, "retry: 3000\n\n"); err != nil {
			return
		}

		if err := writer.Flush(); err != nil {
			return
		}

		for {
			select {
			case payload, open := <-stream:
				if !open {
					return
				}

				if _, err := fmt.Fprintf(writer, "data: %s\n\n", payload); err != nil {
					return
				}
			case <-heartbeat.C:
				if _, err := fmt.Fprint(writer, ": ping\n\n"); err != nil {
					return
				}
			}

			if err := writer.Flush(); err != nil {
				logger.Debug("SSE_CLIENT_GONE", logger.String("tenant.ID", tenantID.String()))

				return
			}
		}
	}))

	return nil
}
