package health

import healthsvc "ava/api/internal/services/health"

type Controller struct {
	healthSvc healthsvc.Service
}

func NewController(healthSvc healthsvc.Service) *Controller {
	return &Controller{
		healthSvc: healthSvc,
	}
}
