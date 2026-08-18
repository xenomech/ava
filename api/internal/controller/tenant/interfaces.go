package tenant

import tenantsvc "ava/api/internal/services/tenant"

type Controller struct {
	tenantService tenantsvc.Service
}

func NewController(tenantService tenantsvc.Service) *Controller {
	return &Controller{tenantService: tenantService}
}
