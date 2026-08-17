package controller

import (
	authctrl "ava/internal/controller/auth"
	healthctrl "ava/internal/controller/health"
	"ava/internal/services"
)

type Controller struct {
	Auth   *authctrl.Controller
	Health *healthctrl.Controller
}

func NewController(service *services.Service) *Controller {
	return &Controller{
		Auth:   authctrl.NewController(service.Auth),
		Health: healthctrl.NewController(service.Health),
	}
}
