package flow

import (
	"context"
	"encoding/json"
	"strings"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	hubrepo "ava/api/internal/repository/hub"
	membershiprepo "ava/api/internal/repository/membership"
	tenantsvc "ava/api/internal/services/tenant"
	"ava/api/pkg/serrors"

	"github.com/google/uuid"
)

type StepHandler interface {
	Validate(ctx context.Context, tenantID, userID uuid.UUID, data json.RawMessage) (model.StepErrors, error)
	Execute(ctx context.Context, tenantID, userID uuid.UUID, data json.RawMessage) error
}

type homeStepHandler struct {
	tenantService  tenantsvc.Service
	membershipRepo membershiprepo.Repository
}

func (h *homeStepHandler) Validate(ctx context.Context, tenantID, userID uuid.UUID, data json.RawMessage) (model.StepErrors, error) {
	if err := requireTenantAdmin(ctx, h.membershipRepo, tenantID, userID); err != nil {
		return nil, err
	}

	payload, err := decodePayload[HomeStepData](data)
	if err != nil {
		return nil, err
	}

	errs := model.StepErrors{}
	if strings.TrimSpace(payload.Name) == "" {
		errs["name"] = "Give your home a name"
	}

	return errs, nil
}

func (h *homeStepHandler) Execute(ctx context.Context, tenantID, userID uuid.UUID, data json.RawMessage) error {
	if err := requireTenantAdmin(ctx, h.membershipRepo, tenantID, userID); err != nil {
		return err
	}

	payload, err := decodePayload[HomeStepData](data)
	if err != nil {
		return err
	}

	_, err = h.tenantService.Update(ctx, tenantID, &dto.UpdateTenantRequest{Name: strings.TrimSpace(payload.Name)})

	return err
}

// hubStepHandler does not pair anything. The client pairs through the hub
// activation endpoint and then submits this step, so all that is left is to
// confirm a hub actually arrived. Someone whose hardware has not shipped yet
// skips the step instead.
type hubStepHandler struct {
	hubRepo hubrepo.Repository
}

func (h *hubStepHandler) Validate(ctx context.Context, tenantID, _ uuid.UUID, _ json.RawMessage) (model.StepErrors, error) {
	hubs, err := h.hubRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	errs := model.StepErrors{}
	if len(hubs) == 0 {
		errs["user_code"] = "Enter the code your hub is showing, or skip for now"
	}

	return errs, nil
}

func (h *hubStepHandler) Execute(context.Context, uuid.UUID, uuid.UUID, json.RawMessage) error {
	return nil
}

func requireTenantAdmin(ctx context.Context, membershipRepo membershiprepo.Repository, tenantID, userID uuid.UUID) error {
	membership, err := membershipRepo.GetByTenantAndUser(ctx, tenantID, userID)
	if err != nil {
		if serrors.Is(err, membershiprepo.ErrMembershipNotFound) {
			return ErrStepNotPermitted
		}

		return err
	}

	if !membership.IsActive() {
		return ErrStepNotPermitted
	}

	if membership.Role != model.TenantRoleOwner && membership.Role != model.TenantRoleAdmin {
		return ErrStepNotPermitted
	}

	return nil
}

func buildHandlers(
	tenantService tenantsvc.Service,
	membershipRepo membershiprepo.Repository,
	hubRepo hubrepo.Repository,
) map[string]StepHandler {
	return map[string]StepHandler{
		"onboarding:home": &homeStepHandler{tenantService: tenantService, membershipRepo: membershipRepo},
		"onboarding:hub":  &hubStepHandler{hubRepo: hubRepo},
	}
}
