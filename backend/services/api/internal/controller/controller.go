package controller

import (
	authctrl "ava/api/internal/controller/auth"
	devicectrl "ava/api/internal/controller/device"
	eventctrl "ava/api/internal/controller/event"
	flowctrl "ava/api/internal/controller/flow"
	healthctrl "ava/api/internal/controller/health"
	hubctrl "ava/api/internal/controller/hub"
	tenantctrl "ava/api/internal/controller/tenant"
	"ava/api/internal/services"
)

type Controller struct {
	Auth   *authctrl.Controller
	Tenant *tenantctrl.Controller
	Flow   *flowctrl.Controller
	Hub    *hubctrl.Controller
	Device *devicectrl.Controller
	Event  *eventctrl.Controller
	Health *healthctrl.Controller
}

func NewController(service *services.Service) *Controller {
	return &Controller{
		Auth:   authctrl.NewController(service.Auth),
		Tenant: tenantctrl.NewController(service.Tenant),
		Flow:   flowctrl.NewController(service.Flow),
		Hub:    hubctrl.NewController(service.Hub),
		Device: devicectrl.NewController(service.Device),
		Event:  eventctrl.NewController(service.Event),
		Health: healthctrl.NewController(service.Health),
	}
}
