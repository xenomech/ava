package tenant

import (
	"context"

	"ava/internal/dto"
	membershiprepo "ava/internal/repository/membership"
	sessionrepo "ava/internal/repository/session"
	tenantrepo "ava/internal/repository/tenant"
	userrepo "ava/internal/repository/user"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, req *dto.CreateTenantRequest) (*dto.TenantResponse, error)
	ListMine(ctx context.Context, userID uuid.UUID) ([]*dto.TenantSummary, error)
	Get(ctx context.Context, tenantID uuid.UUID) (*dto.TenantResponse, error)
	Update(ctx context.Context, tenantID uuid.UUID, req *dto.UpdateTenantRequest) (*dto.TenantResponse, error)
	ListMembers(ctx context.Context, tenantID uuid.UUID) ([]*dto.MemberResponse, error)
	Invite(ctx context.Context, tenantID, invitedByID uuid.UUID, req *dto.InviteMemberRequest) (*dto.MemberResponse, error)
	UpdateMemberRole(ctx context.Context, tenantID, userID uuid.UUID, req *dto.UpdateMemberRoleRequest) error
	RemoveMember(ctx context.Context, tenantID, userID uuid.UUID) error
}

type tenantService struct {
	tenantRepo     tenantrepo.Repository
	membershipRepo membershiprepo.Repository
	userRepo       userrepo.Repository
	sessionRepo    sessionrepo.Repository
}

func NewService(
	tenantRepo tenantrepo.Repository,
	membershipRepo membershiprepo.Repository,
	userRepo userrepo.Repository,
	sessionRepo sessionrepo.Repository,
) Service {
	return &tenantService{
		tenantRepo:     tenantRepo,
		membershipRepo: membershipRepo,
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
	}
}
