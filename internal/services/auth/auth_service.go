package auth

import (
	"context"
	"fmt"
	"time"

	"ava/config"

	"ava/internal/dto"
	"ava/internal/model"
	membershiprepo "ava/internal/repository/membership"
	sessionrepo "ava/internal/repository/session"
	tenantrepo "ava/internal/repository/tenant"
	tokenrepo "ava/internal/repository/token"
	userrepo "ava/internal/repository/user"
	"ava/internal/services/auth/jwt"
	"ava/pkg/email"
	"ava/pkg/logger"
	"ava/pkg/serrors"

	"github.com/google/uuid"
)

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	if !model.IsValidSlug(req.TenantSlug) {
		return nil, ErrInvalidSlug
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		logger.Error("failed to hash password", logger.Err(err))

		return nil, err
	}

	user := model.NewUser(req.Email, req.Username, req.Name, req.Phone, hashedPassword)
	tenant := model.NewTenant(req.TenantName, req.TenantSlug)
	membership := model.NewTenantMembership(tenant.ID, user.ID, model.TenantRoleOwner)

	if err := s.tenantRepo.CreateWithOwner(ctx, user, tenant, membership); err != nil {
		if serrors.Is(err, tenantrepo.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}

		if serrors.Is(err, tenantrepo.ErrTenantAlreadyExists) {
			return nil, ErrTenantAlreadyExists
		}

		logger.Error("failed to create user and tenant", logger.Err(err))

		return nil, err
	}

	token, err := s.generateRandomToken()
	if err != nil {
		logger.Error("failed to generate verification token", logger.Err(err))

		return nil, err
	}

	verificationToken := model.NewToken(
		user.ID,
		token,
		model.TokenTypeEmailVerification,
		time.Now().Add(24*time.Hour),
	)

	if err := s.tokenRepo.CreateToken(ctx, verificationToken); err != nil {
		logger.Error("failed to create verification token", logger.Err(err))

		return nil, err
	}

	emailSvc := email.NewService()
	emailData := map[string]any{
		"Name":            user.Name,
		"VerificationURL": fmt.Sprintf("%s/verify?token=%s", config.GetConfig().AppURL, token),
	}

	if err := emailSvc.Send(ctx, user.Email, "Verify your email address", "verify_email.html", emailData); err != nil {
		logger.Warn("failed to send verification email", logger.String("email", user.Email), logger.Err(err))
	}

	return &dto.RegisterResponse{
		User: *s.userToResponse(user),
		Tenant: dto.TenantSummary{
			ID:   tenant.ID,
			Name: tenant.Name,
			Slug: tenant.Slug,
			Role: membership.Role,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest, deviceInfo dto.DeviceInfo) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		logger.Error("failed to get user by email", logger.Err(err))

		return nil, err
	}

	if !ComparePassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}

	if !user.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	memberships, err := s.activeMemberships(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if len(memberships) == 0 {
		return nil, ErrNoTenantMembership
	}

	membership, err := s.selectMembership(ctx, memberships, req.TenantSlug)
	if err != nil {
		if !serrors.Is(err, ErrTenantSelectionRequired) {
			return nil, err
		}

		tenants, err := s.tenantSummaries(ctx, memberships)
		if err != nil {
			return nil, err
		}

		return &dto.AuthResponse{
			User:                 *s.userToResponse(user),
			NeedsTenantSelection: true,
			Tenants:              tenants,
		}, nil
	}

	return s.issueSession(ctx, user, membership, deviceInfo)
}

func (s *authService) RefreshToken(ctx context.Context, refreshTokenString string) (*dto.TokenResponse, error) {
	claims, err := s.tokenManager.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != jwt.RefreshToken {
		return nil, ErrInvalidToken
	}

	if claims.ID == "" {
		return nil, ErrInvalidToken
	}

	session, err := s.sessionRepo.GetSessionByRID(ctx, claims.TenantID, claims.ID)
	if err != nil {
		if serrors.Is(err, sessionrepo.ErrSessionNotFound) {
			return nil, ErrSessionNotFound
		}

		logger.Error("failed to get session by RID", logger.Err(err))

		return nil, err
	}

	if session.Revoked {
		return nil, ErrSessionRevoked
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	if session.UserID != claims.UserID {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		logger.Error("failed to get user by ID", logger.Err(err))

		return nil, err
	}

	membership, err := s.membershipRepo.GetByTenantAndUser(ctx, claims.TenantID, claims.UserID)
	if err != nil {
		if serrors.Is(err, membershiprepo.ErrMembershipNotFound) {
			return nil, ErrAccessDenied
		}

		logger.Error("failed to get membership", logger.Err(err))

		return nil, err
	}

	if !membership.IsActive() {
		return nil, ErrAccessDenied
	}

	newRID := uuid.NewString()
	if err := s.sessionRepo.UpdateSessionRID(ctx, claims.TenantID, session.ID, newRID); err != nil {
		logger.Error("failed to update session RID", logger.Err(err))

		return nil, err
	}

	return s.issueTokens(user, membership, session.ID, newRID)
}

func (s *authService) ValidateSession(ctx context.Context, tenantID, sessionID uuid.UUID) (*model.Session, error) {
	session, err := s.sessionRepo.GetSessionByID(ctx, tenantID, sessionID)
	if err != nil {
		if serrors.Is(err, sessionrepo.ErrSessionNotFound) {
			return nil, ErrSessionNotFound
		}

		return nil, err
	}

	if session.Revoked {
		return nil, ErrSessionRevoked
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return session, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *authService) Logout(ctx context.Context, tenantID, sessionID uuid.UUID) error {
	if err := s.sessionRepo.RevokeSession(ctx, tenantID, sessionID); err != nil {
		logger.Error("failed to revoke session", logger.Err(err))

		return err
	}

	return nil
}

func (s *authService) CurrentSession(ctx context.Context, tenantID, userID uuid.UUID) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		logger.Error("failed to get user by ID", logger.Err(err))

		return nil, err
	}

	membership, err := s.membershipRepo.GetByTenantAndUser(ctx, tenantID, userID)
	if err != nil {
		if serrors.Is(err, membershiprepo.ErrMembershipNotFound) {
			return nil, ErrAccessDenied
		}

		logger.Error("failed to get membership", logger.Err(err))

		return nil, err
	}

	if !membership.IsActive() {
		return nil, ErrAccessDenied
	}

	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		logger.Error("failed to get tenant by ID", logger.Err(err))

		return nil, err
	}

	return &dto.AuthResponse{
		User: *s.userToResponse(user),
		Tenant: &dto.TenantSummary{
			ID:   tenant.ID,
			Name: tenant.Name,
			Slug: tenant.Slug,
			Role: membership.Role,
		},
	}, nil
}

func (s *authService) LogoutAll(ctx context.Context, tenantID, userID uuid.UUID) error {
	if err := s.sessionRepo.RevokeAllUserSessions(ctx, tenantID, userID); err != nil {
		logger.Error("failed to revoke user sessions", logger.Err(err))

		return err
	}

	return nil
}

func (s *authService) SwitchTenant(ctx context.Context, tenantID, userID, sessionID uuid.UUID, tenantSlug string, deviceInfo dto.DeviceInfo) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		logger.Error("failed to get user by ID", logger.Err(err))

		return nil, err
	}

	memberships, err := s.activeMemberships(ctx, userID)
	if err != nil {
		return nil, err
	}

	membership, err := s.selectMembership(ctx, memberships, tenantSlug)
	if err != nil {
		if serrors.Is(err, ErrTenantSelectionRequired) {
			return nil, ErrAccessDenied
		}

		return nil, err
	}

	if err := s.sessionRepo.RevokeSession(ctx, tenantID, sessionID); err != nil &&
		!serrors.Is(err, sessionrepo.ErrSessionNotFound) {
		logger.Error("failed to revoke session on tenant switch", logger.Err(err))

		return nil, err
	}

	return s.issueSession(ctx, user, membership, deviceInfo)
}

func (s *authService) AcceptInvite(ctx context.Context, inviteToken string) (*dto.TenantSummary, error) {
	membership, err := s.membershipRepo.GetByInviteToken(ctx, inviteToken)
	if err != nil {
		if serrors.Is(err, membershiprepo.ErrMembershipNotFound) {
			return nil, ErrInviteInvalid
		}

		logger.Error("failed to get membership by invite token", logger.Err(err))

		return nil, err
	}

	if err := s.membershipRepo.Activate(ctx, membership.TenantID, membership.UserID); err != nil {
		logger.Error("failed to activate membership", logger.Err(err))

		return nil, err
	}

	tenant, err := s.tenantRepo.GetByID(ctx, membership.TenantID)
	if err != nil {
		logger.Error("failed to get tenant by ID", logger.Err(err))

		return nil, err
	}

	return &dto.TenantSummary{
		ID:   tenant.ID,
		Name: tenant.Name,
		Slug: tenant.Slug,
		Role: membership.Role,
	}, nil
}

func (s *authService) VerifyEmail(ctx context.Context, tokenString string) error {
	token, err := s.tokenRepo.GetValidToken(ctx, tokenString, model.TokenTypeEmailVerification)
	if err != nil {
		if serrors.Is(err, tokenrepo.ErrTokenNotFound) {
			return ErrInvalidToken
		}

		logger.Error("failed to get valid token", logger.Err(err))

		return err
	}

	if err := s.userRepo.VerifyEmail(ctx, token.UserID); err != nil {
		logger.Error("failed to verify email", logger.Err(err))

		return err
	}

	if err := s.tokenRepo.MarkTokenUsed(ctx, token.ID); err != nil {
		logger.Error("failed to mark token as used", logger.Err(err))

		return err
	}

	return nil
}

func (s *authService) ResendVerification(ctx context.Context, emailAddr string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, emailAddr)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return nil
		}

		logger.Error("failed to get user by email", logger.Err(err))

		return err
	}

	if user.EmailVerified {
		return nil
	}

	token, err := s.generateRandomToken()
	if err != nil {
		logger.Error("failed to generate verification token", logger.Err(err))

		return err
	}

	verificationToken := model.NewToken(
		user.ID,
		token,
		model.TokenTypeEmailVerification,
		time.Now().Add(24*time.Hour),
	)

	if err := s.tokenRepo.CreateToken(ctx, verificationToken); err != nil {
		logger.Error("failed to create verification token", logger.Err(err))

		return err
	}

	emailSvc := email.NewService()
	emailData := map[string]any{
		"Name":            user.Name,
		"VerificationURL": fmt.Sprintf("%s/verify?token=%s", config.GetConfig().AppURL, token),
	}

	if err := emailSvc.Send(ctx, user.Email, "Verify your email address", "verify_email.html", emailData); err != nil {
		logger.Warn("failed to send verification email", logger.String("email", user.Email), logger.Err(err))
	}

	return nil
}

func (s *authService) ForgotPassword(ctx context.Context, emailAddr string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, emailAddr)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return nil
		}

		logger.Error("failed to get user by email", logger.Err(err))

		return err
	}

	token, err := s.generateRandomToken()
	if err != nil {
		logger.Error("failed to generate reset token", logger.Err(err))

		return err
	}

	resetToken := model.NewToken(
		user.ID,
		token,
		model.TokenTypePasswordReset,
		time.Now().Add(1*time.Hour),
	)

	if err := s.tokenRepo.CreateToken(ctx, resetToken); err != nil {
		logger.Error("failed to create reset token", logger.Err(err))

		return err
	}

	emailSvc := email.NewService()
	emailData := map[string]any{
		"Name":     user.Name,
		"ResetURL": fmt.Sprintf("%s/reset-password?token=%s", config.GetConfig().AppURL, token),
	}

	if err := emailSvc.Send(ctx, user.Email, "Reset your password", "password_reset.html", emailData); err != nil {
		logger.Warn("failed to send password reset email", logger.String("email", user.Email), logger.Err(err))
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, tokenString, newPassword string) error {
	token, err := s.tokenRepo.GetValidToken(ctx, tokenString, model.TokenTypePasswordReset)
	if err != nil {
		if serrors.Is(err, tokenrepo.ErrTokenNotFound) {
			return ErrInvalidToken
		}

		logger.Error("failed to get valid token", logger.Err(err))

		return err
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		logger.Error("failed to hash password", logger.Err(err))

		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, token.UserID, hashedPassword); err != nil {
		logger.Error("failed to update password", logger.Err(err))

		return err
	}

	if err := s.tokenRepo.MarkTokenUsed(ctx, token.ID); err != nil {
		logger.Error("failed to mark token as used", logger.Err(err))

		return err
	}

	if err := s.sessionRepo.RevokeAllUserSessionsGlobal(ctx, token.UserID); err != nil {
		logger.Error("failed to revoke all sessions after password reset", logger.Err(err))

		return err
	}

	return nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if serrors.Is(err, userrepo.ErrUserNotFound) {
			return ErrUserNotFound
		}

		logger.Error("failed to get user by ID", logger.Err(err))

		return err
	}

	if !ComparePassword(user.Password, oldPassword) {
		return ErrPasswordMismatch
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		logger.Error("failed to hash password", logger.Err(err))

		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		logger.Error("failed to update password", logger.Err(err))

		return err
	}

	return nil
}

func (s *authService) GetUserSessions(ctx context.Context, tenantID, userID, currentSessionID uuid.UUID) ([]*dto.SessionResponse, error) {
	sessions, err := s.sessionRepo.GetUserSessions(ctx, tenantID, userID)
	if err != nil {
		logger.Error("failed to get user sessions", logger.Err(err))

		return nil, err
	}

	responses := make([]*dto.SessionResponse, 0, len(sessions))

	for _, session := range sessions {
		responses = append(responses, &dto.SessionResponse{
			ID:         session.ID,
			DeviceName: session.DeviceName,
			IPAddress:  session.IPAddress,
			UserAgent:  session.UserAgent,
			CreatedAt:  session.CreatedAt,
			ExpiresAt:  session.ExpiresAt,
			Current:    session.ID == currentSessionID,
		})
	}

	return responses, nil
}
