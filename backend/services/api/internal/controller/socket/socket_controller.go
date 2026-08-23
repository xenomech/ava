package socket

import (
	"context"
	"encoding/json"
	"time"

	"ava/api/internal/dto"
	"ava/api/pkg/response"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	pingInterval = 25 * time.Second
	readTimeout  = 70 * time.Second
	writeTimeout = 10 * time.Second
	readLimit    = 4 * 1024
	outboundSize = 16
)

func (c *Controller) Upgrade(ctx *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(ctx) {
		return response.Send(ctx, fiber.StatusUpgradeRequired, nil, "This endpoint expects a websocket upgrade")
	}

	return ctx.Next()
}

func (c *Controller) Serve(conn *websocket.Conn) {
	tenantID, ok := conn.Locals("tenantID").(uuid.UUID)
	if !ok {
		_ = conn.Close()

		return
	}

	stream, unsubscribe := c.events.Subscribe(tenantID)
	defer unsubscribe()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	outbound := make(chan []byte, outboundSize)

	go c.write(ctx, conn, stream, outbound)

	c.read(ctx, conn, tenantID, outbound)

	logger.Debug("SOCKET_CLOSED", logger.Any("tenant.ID", tenantID))
}

func (c *Controller) read(ctx context.Context, conn *websocket.Conn, tenantID uuid.UUID, outbound chan<- []byte) {
	conn.SetReadLimit(readLimit)
	_ = conn.SetReadDeadline(time.Now().Add(readTimeout))

	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(readTimeout))
	})

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			return
		}

		_ = conn.SetReadDeadline(time.Now().Add(readTimeout))

		var frame dto.DeviceCommandFrame

		if err := json.Unmarshal(payload, &frame); err != nil || frame.Type != dto.EventDeviceCommand {
			continue
		}

		c.dispatch(ctx, tenantID, &frame, outbound)
	}
}

func (c *Controller) dispatch(
	ctx context.Context,
	tenantID uuid.UUID,
	frame *dto.DeviceCommandFrame,
	outbound chan<- []byte,
) {
	req := &dto.SendCommandRequest{Trait: frame.Trait, Value: frame.Value}

	if _, err := c.devices.SendCommand(ctx, tenantID, frame.DeviceID, req); err != nil {
		reject(outbound, frame.DeviceID, err)
	}
}

func reject(outbound chan<- []byte, deviceID uuid.UUID, err error) {
	code := serrors.CodeOf(err)
	if code == "" {
		code = "command_failed"
	}

	payload, marshalErr := json.Marshal(dto.NewCommandRejectedEvent(deviceID, code, err.Error()))
	if marshalErr != nil {
		return
	}

	select {
	case outbound <- payload:
	default:
	}
}

func (c *Controller) write(ctx context.Context, conn *websocket.Conn, stream, outbound <-chan []byte) {
	ping := time.NewTicker(pingInterval)

	defer func() {
		ping.Stop()
		_ = conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case payload, open := <-stream:
			if !open || !send(conn, websocket.TextMessage, payload) {
				return
			}
		case payload := <-outbound:
			if !send(conn, websocket.TextMessage, payload) {
				return
			}
		case <-ping.C:
			if !send(conn, websocket.PingMessage, nil) {
				return
			}
		}
	}
}

func send(conn *websocket.Conn, kind int, payload []byte) bool {
	if err := conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		return false
	}

	return conn.WriteMessage(kind, payload) == nil
}
