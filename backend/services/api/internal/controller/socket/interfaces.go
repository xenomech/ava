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
	life    context.Context
	events  eventsvc.Service
	devices Commander
}

func NewController(life context.Context, events eventsvc.Service, devices devicesvc.Service) *Controller {
	return &Controller{life: life, events: events, devices: devices}
}
