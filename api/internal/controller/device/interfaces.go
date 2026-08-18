package device

import devicesvc "ava/api/internal/services/device"

type Controller struct {
	deviceService devicesvc.Service
}

func NewController(deviceService devicesvc.Service) *Controller {
	return &Controller{deviceService: deviceService}
}
