package flow

import (
	"context"
	"encoding/json"
	"strings"

	"ava/api/internal/dto"
	"ava/api/internal/model"
	membershiprepo "ava/api/internal/repository/membership"
	userrepo "ava/api/internal/repository/user"
	tenantsvc "ava/api/internal/services/tenant"
	"ava/api/pkg/serrors"
	"ava/pkg/logger"

	"github.com/google/uuid"
)

type StepHandler interface {
	Validate(ctx context.Context, tenantID, userID uuid.UUID, data json.RawMessage) (model.StepErrors, error)
	Execute(ctx context.Context, tenantID, userID uuid.UUID, data json.RawMessage) error
}

type profileStepHandler struct {
	userRepo userrepo.Repository
}

func (h *profileStepHandler) Validate(_ context.Context, _, _ uuid.UUID, data json.RawMessage) (model.StepErrors, error) {
	payload, err := decodePayload[ProfileStepData](data)
	if err != nil {
		return nil, err
	}

	errs := model.StepErrors{}
	if strings.TrimSpace(payload.Name) == "" {
		errs["name"] = "Your name is required"
	}

	return errs, nil
}

func (h *profileStepHandler) Execute(ctx context.Context, _, userID uuid.UUID, data json.RawMessage) error {
	payload, err := decodePayload[ProfileStepData](data)
	if err != nil {
		return err
	}

	update := &model.User{
		Name:  strings.TrimSpace(payload.Name),
		Phone: strings.TrimSpace(payload.Phone),
	}
	update.ID = userID

	return h.userRepo.UpdateUser(ctx, update)
}

type workspaceStepHandler struct {
	tenantService  tenantsvc.Service
	membershipRepo membershiprepo.Repository
}

func (h *workspaceStepHandler) Validate(ctx context.Context, tenantID, userID uuid.UUID, data json.RawMessage) (model.StepErrors, error) {
	if err := requireTenantAdmin(ctx, h.membershipRepo, tenantID, userID); err != nil {
		return nil, err
	}

	payload, err := decodePayload[WorkspaceStepData](data)
	if err != nil {
		return nil, err
	}

	errs := model.StepErrors{}
	if strings.TrimSpace(payload.Name) == "" {
		errs["name"] = "Workspace name is required"
	}

	return errs, nil
}

func (h *workspaceStepHandler) Execute(ctx context.Context, tenantID, userID uuid.UUID, data json.RawMessage) error {
	if err := requireTenantAdmin(ctx, h.membershipRepo, tenantID, userID); err != nil {
		return err
	}

	payload, err := decodePayload[WorkspaceStepData](data)
	if err != nil {
		return err
	}

	_, err = h.tenantService.Update(ctx, tenantID, &dto.UpdateTenantRequest{Name: strings.TrimSpace(payload.Name)})

	return err
}

type inviteTeamStepHandler struct {
	tenantService  tenantsvc.Service
	userRepo       userrepo.Repository
	membershipRepo membershiprepo.Repository
}

func (h *inviteTeamStepHandler) Validate(ctx context.Context, tenantID, userID uuid.UUID, data json.RawMessage) (model.StepErrors, error) {
	if err := requireTenantAdmin(ctx, h.membershipRepo, tenantID, userID); err != nil {
		return nil, err
	}

	payload, err := decodePayload[InviteTeamStepData](data)
	if err != nil {
		return nil, err
	}

	errs := model.StepErrors{}

	emails := cleanEmails(payload.Emails)
	if len(emails) == 0 {
		errs["emails"] = "Add at least one email address, or skip this step"

		return errs, nil
	}

	if payload.Role != "" && (!payload.Role.IsValid() || payload.Role == model.TenantRoleOwner) {
		errs["role"] = "Choose either admin or member"
	}

	for _, addr := range emails {
		if _, err := h.userRepo.GetUserByEmail(ctx, addr); err != nil {
			if serrors.Is(err, userrepo.ErrUserNotFound) {
				errs["emails"] = "No account exists for " + addr + " — they need to sign up first"

				break
			}

			return nil, err
		}
	}

	return errs, nil
}

func (h *inviteTeamStepHandler) Execute(ctx context.Context, tenantID, userID uuid.UUID, data json.RawMessage) error {
	if err := requireTenantAdmin(ctx, h.membershipRepo, tenantID, userID); err != nil {
		return err
	}

	payload, err := decodePayload[InviteTeamStepData](data)
	if err != nil {
		return err
	}

	role := model.TenantRoleMember
	if payload.Role != "" {
		role = payload.Role
	}

	for _, addr := range cleanEmails(payload.Emails) {
		_, err := h.tenantService.Invite(ctx, tenantID, userID, &dto.InviteMemberRequest{Email: addr, Role: role})
		if err == nil {
			continue
		}

		if serrors.Is(err, tenantsvc.ErrAlreadyMember) || serrors.Is(err, tenantsvc.ErrAlreadyInvited) {
			logger.Info("flow.invite_team.skipped", logger.String("email", addr), logger.Err(err))

			continue
		}

		return err
	}

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

func cleanEmails(emails []string) []string {
	cleaned := make([]string, 0, len(emails))

	for _, addr := range emails {
		if trimmed := strings.TrimSpace(addr); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	return cleaned
}

func buildHandlers(
	tenantService tenantsvc.Service,
	userRepo userrepo.Repository,
	membershipRepo membershiprepo.Repository,
) map[string]StepHandler {
	return map[string]StepHandler{
		"onboarding:profile":   &profileStepHandler{userRepo: userRepo},
		"onboarding:workspace": &workspaceStepHandler{tenantService: tenantService, membershipRepo: membershipRepo},
		"onboarding:invite_team": &inviteTeamStepHandler{
			tenantService:  tenantService,
			userRepo:       userRepo,
			membershipRepo: membershipRepo,
		},
	}
}
