package controller

import (
	authctrl "ava/internal/controller/auth"
	flowctrl "ava/internal/controller/flow"
	healthctrl "ava/internal/controller/health"
	tenantctrl "ava/internal/controller/tenant"
	"ava/internal/services"
)

type Controller struct {
	Auth   *authctrl.Controller
	Tenant *tenantctrl.Controller
	Flow   *flowctrl.Controller
	Health *healthctrl.Controller
}

func NewController(service *services.Service) *Controller {
	return &Controller{
		Auth:   authctrl.NewController(service.Auth),
		Tenant: tenantctrl.NewController(service.Tenant),
		Flow:   flowctrl.NewController(service.Flow),
		Health: healthctrl.NewController(service.Health),
	}
}
