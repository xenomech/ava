package controller

import (
	healthctrl "ava/internal/controller/health"
	"ava/internal/services"
)

type Controller struct {
	Health *healthctrl.Controller
}

func NewController(service *services.Service) *Controller {
	return &Controller{
		Health: healthctrl.NewController(service.Health),
	}
}
