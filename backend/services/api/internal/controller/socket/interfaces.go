package socket

import (
	"context"

	"ava/api/internal/dto"
	devicesvc "ava/api/internal/services/device"
	eventsvc "ava/api/internal/services/event"

	"github.com/google/uuid"
)

type Commander interface {
	SendCommand(
		ctx context.Context,
		tenantID, deviceID uuid.UUID,
		req *dto.SendCommandRequest,
	) (*dto.CommandAcceptedResponse, error)
}

type Controller struct {
	events  eventsvc.Service
	devices Commander
}

func NewController(events eventsvc.Service, devices devicesvc.Service) *Controller {
	return &Controller{events: events, devices: devices}
}
