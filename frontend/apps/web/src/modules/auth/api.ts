import {
  authResponseDto,
  authenticatedDto,
  registerResponseDto,
  sessionListDto,
  tokensDto,
  type AuthResponseDto,
  type LoginRequest,
  type RegisterRequest,
  type RegisterResponseDto,
  type SessionDto,
} from "@ava/contracts";
import type { z } from "zod";

import { request } from "@/config/http/request";

const sessionDtoShape = authenticatedDto.omit({ tokens: true }).extend({
  tokens: tokensDto.optional(),
});

export type Session = z.infer<typeof sessionDtoShape>;

export function fetchSession(signal?: AbortSignal): Promise<Session> {
  return request({ url: "/auth/me", schema: sessionDtoShape, signal });
}

export function login(body: LoginRequest, device?: string): Promise<AuthResponseDto> {
  return request({
    method: "post",
    url: "/auth/login",
    body,
    schema: authResponseDto,
    headers: device ? { "X-Device-Name": device } : undefined,
  });
}

export function register(body: RegisterRequest): Promise<RegisterResponseDto> {
  return request({ method: "post", url: "/auth/register", body, schema: registerResponseDto });
}

export function switchTenant(tenantSlug: string): Promise<AuthResponseDto> {
  return request({
    method: "post",
    url: "/auth/switch-tenant",
    body: { tenant_slug: tenantSlug },
    schema: authResponseDto,
  });
}

export function logout(): Promise<void> {
  return request({ method: "post", url: "/auth/logout" });
}

export function listSessions(signal?: AbortSignal): Promise<SessionDto[]> {
  return request({ url: "/auth/sessions", schema: sessionListDto, signal });
}

export function verifyEmail(token: string): Promise<void> {
  return request({ method: "post", url: "/auth/verify-email", body: { token } });
}

export function forgotPassword(email: string): Promise<void> {
  return request({ method: "post", url: "/auth/forgot-password", body: { email } });
}

export function resetPassword(token: string, newPassword: string): Promise<void> {
  return request({
    method: "post",
    url: "/auth/reset-password",
    body: { token, new_password: newPassword },
  });
}

export function acceptInvite(token: string): Promise<void> {
  return request({ method: "post", url: "/auth/accept-invite", body: { token } });
}
