package hub

import (
	"context"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	hubrepo "ava/api/internal/repository/hub"
	tenantrepo "ava/api/internal/repository/tenant"
	"ava/api/internal/services/auth/jwt"
	eventsvc "ava/api/internal/services/event"

	"github.com/google/uuid"
)

type Service interface {
	RequestCode(ctx context.Context, req *dto.DeviceCodeRequest) (*dto.DeviceCodeResponse, error)
	Poll(ctx context.Context, deviceCode string) (*dto.HubTokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.HubTokenResponse, error)
	Activate(ctx context.Context, tenantID, userID uuid.UUID, userCode string) (*dto.HubResponse, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*dto.HubResponse, error)
	Revoke(ctx context.Context, tenantID, hubID uuid.UUID) error
	ApplyPresence(ctx context.Context, hubID uuid.UUID, online bool) (*model.Hub, error)
	ValidateDevice(ctx context.Context, hubID uuid.UUID) (*model.Hub, error)
}

type hubService struct {
	hubRepo      hubrepo.Repository
	tenantRepo   tenantrepo.Repository
	events       eventsvc.Service
	tokenManager jwt.TokenManager
}

func NewService(hubRepo hubrepo.Repository, tenantRepo tenantrepo.Repository, events eventsvc.Service) Service {
	return &hubService{
		hubRepo:      hubRepo,
		tenantRepo:   tenantRepo,
		events:       events,
		tokenManager: jwt.NewTokenManager(),
	}
}
