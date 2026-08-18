package controller

import (
	authctrl "ava/api/internal/controller/auth"
	flowctrl "ava/api/internal/controller/flow"
	healthctrl "ava/api/internal/controller/health"
	tenantctrl "ava/api/internal/controller/tenant"
	"ava/api/internal/services"
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
