package middleware

import "ava/internal/services"

type Middleware struct {
}

func NewMiddleware(service *services.Service) *Middleware {
	return &Middleware{}
}
