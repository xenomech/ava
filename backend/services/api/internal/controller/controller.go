package controller

import (
	apitokenctrl "ava/api/internal/controller/apitoken"
	authctrl "ava/api/internal/controller/auth"
	devicectrl "ava/api/internal/controller/device"
	flowctrl "ava/api/internal/controller/flow"
	healthctrl "ava/api/internal/controller/health"
	hubctrl "ava/api/internal/controller/hub"
	roomctrl "ava/api/internal/controller/room"
	scenectrl "ava/api/internal/controller/scene"
	socketctrl "ava/api/internal/controller/socket"
	tenantctrl "ava/api/internal/controller/tenant"
	"ava/api/internal/services"
)

type Controller struct {
	APIToken *apitokenctrl.Controller
	Auth     *authctrl.Controller
	Tenant   *tenantctrl.Controller
	Flow     *flowctrl.Controller
	Room     *roomctrl.Controller
	Scene    *scenectrl.Controller
	Hub      *hubctrl.Controller
	Device   *devicectrl.Controller
	Socket   *socketctrl.Controller
	Health   *healthctrl.Controller
}

func NewController(service *services.Service) *Controller {
	return &Controller{
		APIToken: apitokenctrl.NewController(service.APIToken),
		Auth:     authctrl.NewController(service.Auth),
		Tenant:   tenantctrl.NewController(service.Tenant),
		Flow:     flowctrl.NewController(service.Flow),
		Room:     roomctrl.NewController(service.Room),
		Scene:    scenectrl.NewController(service.Scene),
		Hub:      hubctrl.NewController(service.Hub),
		Device:   devicectrl.NewController(service.Device),
		Socket:   socketctrl.NewController(service.Event, service.Device),
		Health:   healthctrl.NewController(service.Health),
	}
}
