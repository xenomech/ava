package tenant

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"ava/api/internal/dto"
	"ava/api/internal/model"
)

func generateInviteToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate invite token: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

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
