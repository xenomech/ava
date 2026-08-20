package hub

import hubsvc "ava/api/internal/services/hub"

type Controller struct {
	hubService hubsvc.Service
}

func NewController(hubService hubsvc.Service) *Controller {
	return &Controller{hubService: hubService}
}
