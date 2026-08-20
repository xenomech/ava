package auth

import authsvc "ava/api/internal/services/auth"

type Controller struct {
	authService authsvc.Service
}

func NewController(authService authsvc.Service) *Controller {
	return &Controller{authService: authService}
}
