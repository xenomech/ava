package tenant

import (
	"context"
	"fmt"
	"time"

	"ava/api/config"
	"ava/api/internal/dto"
	"ava/api/internal/model"
	membershiprepo "ava/api/internal/repository/membership"
	tenantrepo "ava/api/internal/repository/tenant"
	userrepo "ava/api/internal/repository/user"
	"ava/api/pkg/email"
	"ava/api/pkg/logger"
	"ava/api/pkg/serrors"

	"github.com/google/uuid"
)

const inviteExpiry = 7 * 24 * time.Hour

func (s *tenantService) Create(ctx context.Context, userID uuid.UUID, req *dto.CreateTenantRequest) (*dto.TenantResponse, error) {
	if !model.IsValidSlug(req.Slug) {
		return nil, ErrInvalidSlug
	}

	newTenant := model.NewTenant(req.Name, req.Slug)
	membership := model.NewTenantMembership(newTenant.ID, userID, model.TenantRoleOwner)

	if err := s.tenantRepo.CreateWithMembership(ctx, newTenant, membership); err != nil {
		if serrors.Is(err, tenantrepo.ErrTenantAlreadyExists) {
			return nil, ErrTenantAlreadyExists
		}

		logger.Error("tenant.Create", logger.Err(err))

		return nil, err
	}

	return toTenantResponse(newTenant), nil
}

func (s *tenantService) ListMine(ctx context.Context, userID uuid.UUID) ([]*dto.TenantSummary, error) {
	memberships, err := s.membershipRepo.ListByUser(ctx, userID)
	if err != nil {
		logger.Error("tenant.ListMine", logger.Err(err))

		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(memberships))
	roles := make(map[uuid.UUID]model.TenantRole, len(memberships))

	for _, membership := range memberships {
		if !membership.IsActive() {
			continue
		}

		ids = append(ids, membership.TenantID)
		roles[membership.TenantID] = membership.Role
	}

	tenants, err := s.tenantRepo.ListByIDs(ctx, ids)
	if err != nil {
		logger.Error("tenant.ListMine", logger.Err(err))

		return nil, err
	}

	summaries := make([]*dto.TenantSummary, 0, len(tenants))

	for _, listed := range tenants {
		summaries = append(summaries, &dto.TenantSummary{
			ID:   listed.ID,
			Name: listed.Name,
			Slug: listed.Slug,
			Role: roles[listed.ID],
		})
	}

	return summaries, nil
}

func (s *tenantService) Get(ctx context.Context, tenantID uuid.UUID) (*dto.TenantResponse, error) {
	found, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		if serrors.Is(err, tenantrepo.ErrTenantNotFound) {
			return nil, ErrTenantNotFound
		}

		logger.Error("tenant.Get", logger.Err(err))

		return nil, err
	}

	return toTenantResponse(found), nil
}

func (s *tenantService) Update(ctx context.Context, tenantID uuid.UUID, req *dto.UpdateTenantRequest) (*dto.TenantResponse, error) {
	found, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		if serrors.Is(err, tenantrepo.ErrTenantNotFound) {
			return nil, ErrTenantNotFound
		}

		logger.Error("tenant.Update", logger.Err(err))

		return nil, err
	}

	found.Name = req.Name

	if err := s.tenantRepo.Update(ctx, found); err != nil {
		if serrors.Is(err, tenantrepo.ErrTenantNotFound) {
			return nil, ErrTenantNotFound
		}

		logger.Error("tenant.Update", logger.Err(err))

		return nil, err
	}

	return toTenantResponse(found), nil
}

func (s *tenantService) ListMembers(ctx context.Context, tenantID uuid.UUID) ([]*dto.MemberResponse, error) {
	memberships, err := s.membershipRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		logger.Error("tenant.ListMembers", logger.Err(err))

		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(memberships))
	for _, membership := range memberships {
		ids = append(ids, membership.UserID)
	}

	users, err := s.userRepo.ListByIDs(ctx, ids)
	if err != nil {
		logger.Error("tenant.ListMembers", logger.Err(err))

		return nil, err
	}

	byID := make(map[uuid.UUID]*model.User, len(users))
	for _, user := range users {
		byID[user.ID] = user
	}

	responses := make([]*dto.MemberResponse, 0, len(memberships))
	for _, membership := range memberships {
		responses = append(responses, toMemberResponse(membership, byID[membership.UserID]))
	}

	return responses, nil
}

func (s *tenantService) Invite(ctx context.Context, tenantID, invitedByID uuid.UUID, req *dto.InviteMemberRequest) (*dto.MemberResponse, error) {
	if !req.Role.IsValid() || req.Role == model.TenantRoleOwner {
		return nil, ErrInvalidRole
	}

	invitee, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		logger.Error("tenant.Invite", logger.Err(err))

		return nil, err
	}

	existing, err := s.membershipRepo.GetByTenantAndUser(ctx, tenantID, invitee.ID)
	if err != nil && !serrors.Is(err, membershiprepo.ErrMembershipNotFound) {
		logger.Error("tenant.Invite", logger.Err(err))

		return nil, err
	}

	if existing != nil {
		if existing.IsActive() {
			return nil, ErrAlreadyMember
		}

		return nil, ErrAlreadyInvited
	}

	inviteToken, err := generateInviteToken()
	if err != nil {
		logger.Error("tenant.Invite", logger.Err(err))

		return nil, err
	}

	membership := model.NewTenantInvite(tenantID, invitee.ID, invitedByID, req.Role, inviteToken, time.Now().Add(inviteExpiry))

	if err := s.membershipRepo.Create(ctx, membership); err != nil {
		if serrors.Is(err, membershiprepo.ErrMembershipAlreadyExists) {
			return nil, ErrAlreadyMember
		}

		logger.Error("tenant.Invite", logger.Err(err))

		return nil, err
	}

	invitingTenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		logger.Error("tenant.Invite", logger.Err(err))

		return nil, err
	}

	emailSvc := email.NewService()
	emailData := map[string]any{
		"Name":       invitee.Name,
		"TenantName": invitingTenant.Name,
		"Role":       string(req.Role),
		"InviteURL":  fmt.Sprintf("%s/accept-invite?token=%s", config.GetConfig().AppURL, inviteToken),
	}

	if err := emailSvc.Send(ctx, invitee.Email, "You have been invited to "+invitingTenant.Name, "tenant_invite.html", emailData); err != nil {
		logger.Warn("failed to send invitation email", logger.String("email", invitee.Email), logger.Err(err))
	}

	return toMemberResponse(membership, invitee), nil
}

func (s *tenantService) UpdateMemberRole(ctx context.Context, tenantID, userID uuid.UUID, req *dto.UpdateMemberRoleRequest) error {
	if !req.Role.IsValid() {
		return ErrInvalidRole
	}

	membership, err := s.membershipRepo.GetByTenantAndUser(ctx, tenantID, userID)
	if err != nil {
		if serrors.Is(err, membershiprepo.ErrMembershipNotFound) {
			return ErrMemberNotFound
		}

		logger.Error("tenant.UpdateMemberRole", logger.Err(err))

		return err
	}

	if membership.Role == model.TenantRoleOwner && req.Role != model.TenantRoleOwner {
		lastOwner, err := s.isLastOwner(ctx, tenantID)
		if err != nil {
			return err
		}

		if lastOwner {
			return ErrLastOwner
		}
	}

	if err := s.membershipRepo.UpdateRole(ctx, tenantID, userID, req.Role); err != nil {
		if serrors.Is(err, membershiprepo.ErrMembershipNotFound) {
			return ErrMemberNotFound
		}

		logger.Error("tenant.UpdateMemberRole", logger.Err(err))

		return err
	}

	if err := s.sessionRepo.RevokeAllUserSessions(ctx, tenantID, userID); err != nil {
		logger.Error("tenant.UpdateMemberRole", logger.Err(err))

		return err
	}

	return nil
}

func (s *tenantService) RemoveMember(ctx context.Context, tenantID, userID uuid.UUID) error {
	membership, err := s.membershipRepo.GetByTenantAndUser(ctx, tenantID, userID)
	if err != nil {
		if serrors.Is(err, membershiprepo.ErrMembershipNotFound) {
			return ErrMemberNotFound
		}

		logger.Error("tenant.RemoveMember", logger.Err(err))

		return err
	}

	if membership.Role == model.TenantRoleOwner {
		lastOwner, err := s.isLastOwner(ctx, tenantID)
		if err != nil {
			return err
		}

		if lastOwner {
			return ErrLastOwner
		}
	}

	if err := s.membershipRepo.Delete(ctx, tenantID, userID); err != nil {
		if serrors.Is(err, membershiprepo.ErrMembershipNotFound) {
			return ErrMemberNotFound
		}

		logger.Error("tenant.RemoveMember", logger.Err(err))

		return err
	}

	if err := s.sessionRepo.RevokeAllUserSessions(ctx, tenantID, userID); err != nil {
		logger.Error("tenant.RemoveMember", logger.Err(err))

		return err
	}

	return nil
}

func (s *tenantService) isLastOwner(ctx context.Context, tenantID uuid.UUID) (bool, error) {
	owners, err := s.membershipRepo.CountByRole(ctx, tenantID, model.TenantRoleOwner)
	if err != nil {
		logger.Error("tenant.isLastOwner", logger.Err(err))

		return false, err
	}

	return owners <= 1, nil
}
