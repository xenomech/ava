import { z } from "zod";

import { tenantSlugSchema } from "../tenant/tenant.request";

export const passwordSchema = z.string().min(8, "Must be at least 8 characters");

export const loginRequest = z.object({
  email: z.email(),
  password: z.string().min(1, "Password is required"),
  tenant_slug: z.string().optional(),
});

export type LoginRequest = z.infer<typeof loginRequest>;

export const registerRequest = z.object({
  email: z.email(),
  username: z.string().min(3).max(30),
  name: z.string().min(1, "Name is required"),
  password: passwordSchema,
  phone: z.string().optional(),
  tenant_name: z.string().min(1, "Workspace name is required"),
  tenant_slug: tenantSlugSchema,
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
