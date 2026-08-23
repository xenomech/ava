import { z } from "zod";

export const passwordSchema = z.string().min(8, "Must be at least 8 characters");

export const loginRequest = z.object({
  email: z.email(),
  password: z.string().min(1, "Password is required"),
  tenant_slug: z.string().optional(),
});

export type LoginRequest = z.infer<typeof loginRequest>;

// Three fields. A username and a home slug are machine identifiers the server
// derives; the home is named during onboarding, not here.
export const registerRequest = z.object({
  name: z.string().min(1, "Name is required"),
  email: z.email(),
  password: passwordSchema,
});

export type RegisterRequest = z.infer<typeof registerRequest>;

export const refreshTokenRequest = z.object({
  refresh_token: z.string().min(1),
});

export type RefreshTokenRequest = z.infer<typeof refreshTokenRequest>;

export const switchTenantRequest = z.object({
  tenant_slug: z.string().min(1),
});

export type SwitchTenantRequest = z.infer<typeof switchTenantRequest>;

export const verifyEmailRequest = z.object({ token: z.string().min(1) });
export const acceptInviteRequest = z.object({ token: z.string().min(1) });
export const resendVerificationRequest = z.object({ email: z.email() });
export const forgotPasswordRequest = z.object({ email: z.email() });

export const resetPasswordRequest = z.object({
  token: z.string().min(1),
  new_password: passwordSchema,
});

export type ResetPasswordRequest = z.infer<typeof resetPasswordRequest>;

export const changePasswordRequest = z.object({
  old_password: z.string().min(1, "Current password is required"),
  new_password: passwordSchema,
});

export type ChangePasswordRequest = z.infer<typeof changePasswordRequest>;
