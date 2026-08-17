package services

import (
	"ava/internal/repository"
	healthsvc "ava/internal/services/health"
)

type Service struct {
	Health healthsvc.Service
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Health: healthsvc.NewService(),
	}
}
