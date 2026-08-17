package tenant

import (
	"ava/internal/dto"
	"ava/internal/model"
)

func toTenantResponse(tenant *model.Tenant) *dto.TenantResponse {
	return &dto.TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Slug:      tenant.Slug,
		CreatedAt: tenant.CreatedAt,
	}
}

func toMemberResponse(membership *model.TenantMembership, user *model.User) *dto.MemberResponse {
	response := &dto.MemberResponse{
		UserID:   membership.UserID,
		Role:     membership.Role,
		Status:   membership.Status,
		JoinedAt: membership.JoinedAt,
	}

	if user != nil {
		response.Email = user.Email
		response.Name = user.Name
	}

	return response
}
