package apitoken

import apitokensvc "ava/api/internal/services/apitoken"

type Controller struct {
	apiTokenService apitokensvc.Service
}

func NewController(apiTokenService apitokensvc.Service) *Controller {
	return &Controller{apiTokenService: apiTokenService}
}
