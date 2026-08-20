import { z } from "zod";

export const GENERIC_ERROR_CODES = [
  "bad_request",
  "unauthorized",
  "forbidden",
  "not_found",
  "conflict",
  "validation_error",
  "rate_limited",
  "auth_rate_limited",
  "internal_error",
] as const;

export const DOMAIN_ERROR_CODES = [
  "invalid_credentials",
  "email_not_verified",
  "password_mismatch",
  "session_expired",
  "session_revoked",
  "session_not_found",
  "user_already_exists",
  "user_not_found",
  "invalid_token",
  "access_denied",
  "no_tenant_membership",
  "tenant_selection_required",
  "invite_invalid",
  "hub_revoked",
  "invalid_refresh_token",
  "authorization_pending",
  "slow_down",
  "expired_token",
  "tenant_not_found",
  "tenant_already_exists",
  "invalid_slug",
  "invalid_role",
  "member_not_found",
  "already_member",
  "already_invited",
  "last_owner",
  "invalid_flow_type",
  "flow_already_completed",
  "step_not_found",
  "step_not_current",
  "step_not_skippable",
  "no_previous_step",
  "step_validation_failed",
  "invalid_step_data",
  "step_not_permitted",
  "device_not_found",
  "hub_offline",
  "command_channel_unavailable",
  "invalid_code",
] as const;

export const ERROR_CODES = [...GENERIC_ERROR_CODES, ...DOMAIN_ERROR_CODES] as const;

export type GenericErrorCode = (typeof GENERIC_ERROR_CODES)[number];
export type DomainErrorCode = (typeof DOMAIN_ERROR_CODES)[number];
export type ErrorCode = (typeof ERROR_CODES)[number];

export const errorCodeSchema = z.union([z.enum(ERROR_CODES), z.string()]);

export const errorBodySchema = z.object({
  code: errorCodeSchema,
  message: z.string(),
  details: z.record(z.string(), z.string()).optional(),
});

export type ErrorBody = z.infer<typeof errorBodySchema>;
