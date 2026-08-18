package flow

import flowsvc "ava/api/internal/services/flow"

type Controller struct {
	flowService flowsvc.Service
}

func NewController(flowService flowsvc.Service) *Controller {
	return &Controller{flowService: flowService}
}
