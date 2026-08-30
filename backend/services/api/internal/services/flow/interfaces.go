package flow

import (
	"context"

	"ava/api/internal/dto"
	flowrepo "ava/api/internal/repository/flow"
	hubrepo "ava/api/internal/repository/hub"
	membershiprepo "ava/api/internal/repository/membership"
	tenantsvc "ava/api/internal/services/tenant"

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
	membershipRepo membershiprepo.Repository,
	hubRepo hubrepo.Repository,
) Service {
	return &flowService{
		flowRepo: flowRepo,
		handlers: buildHandlers(tenantService, membershipRepo, hubRepo),
	}
}
