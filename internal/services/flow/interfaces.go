package flow

import (
	"context"

	"ava/internal/dto"
	flowrepo "ava/internal/repository/flow"
	membershiprepo "ava/internal/repository/membership"
	userrepo "ava/internal/repository/user"
	tenantsvc "ava/internal/services/tenant"

	"github.com/google/uuid"
)

type Service interface {
	GetFlow(ctx context.Context, tenantID, userID uuid.UUID, flowType string) (*dto.FlowStateResponse, error)
	SubmitStep(ctx context.Context, tenantID, userID uuid.UUID, flowType, stepID string, req *dto.SubmitStepRequest) (*dto.FlowStateResponse, error)
	GoBack(ctx context.Context, tenantID, userID uuid.UUID, flowType string) (*dto.FlowStateResponse, error)
	SkipStep(ctx context.Context, tenantID, userID uuid.UUID, flowType string) (*dto.FlowStateResponse, error)
}

type flowService struct {
	flowRepo flowrepo.Repository
	handlers map[string]StepHandler
}

func NewService(
	flowRepo flowrepo.Repository,
	tenantService tenantsvc.Service,
	userRepo userrepo.Repository,
	membershipRepo membershiprepo.Repository,
) Service {
	return &flowService{
		flowRepo: flowRepo,
		handlers: buildHandlers(tenantService, userRepo, membershipRepo),
	}
}
