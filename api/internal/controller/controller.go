package controller

import (
	authctrl "ava/api/internal/controller/auth"
	devicectrl "ava/api/internal/controller/device"
	flowctrl "ava/api/internal/controller/flow"
	healthctrl "ava/api/internal/controller/health"
	tenantctrl "ava/api/internal/controller/tenant"
	"ava/api/internal/services"
)

type Controller struct {
	Auth   *authctrl.Controller
	Tenant *tenantctrl.Controller
	Flow   *flowctrl.Controller
	Device *devicectrl.Controller
	Health *healthctrl.Controller
}

func NewController(service *services.Service) *Controller {
	return &Controller{
		Auth:   authctrl.NewController(service.Auth),
		Tenant: tenantctrl.NewController(service.Tenant),
		Flow:   flowctrl.NewController(service.Flow),
		Device: devicectrl.NewController(service.Device),
		Health: healthctrl.NewController(service.Health),
	}
}
