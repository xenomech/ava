package auth

import (
	"ava/internal/dto"
	"ava/internal/model"
)

func (s *authService) userToResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		Name:            user.Name,
		Phone:           user.Phone,
		EmailVerified:   user.EmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:       user.CreatedAt,
	}
}
