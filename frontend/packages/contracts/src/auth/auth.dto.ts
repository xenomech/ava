import { z } from "zod";

import { tenantSummaryDto } from "../tenant/tenant.dto";
import { userDto } from "../user/user.dto";

export const tokensDto = z.object({
  access_token: z.string(),
  refresh_token: z.string(),
    expires_in: z.number().int(),
});

export type TokensDto = z.infer<typeof tokensDto>;

export const authenticatedDto = z.object({
  user: userDto,
  tenant: tenantSummaryDto,
  tokens: tokensDto,
  needs_tenant_selection: z.literal(false).optional(),
});

export const tenantSelectionRequiredDto = z.object({
  user: userDto,
  needs_tenant_selection: z.literal(true),
  tenants: z.array(tenantSummaryDto),
});

export const authResponseDto = z.union([tenantSelectionRequiredDto, authenticatedDto]);

export type AuthenticatedDto = z.infer<typeof authenticatedDto>;
export type TenantSelectionRequiredDto = z.infer<typeof tenantSelectionRequiredDto>;
export type AuthResponseDto = z.infer<typeof authResponseDto>;

export function needsTenantSelection(
  response: AuthResponseDto,
): response is TenantSelectionRequiredDto {
  return response.needs_tenant_selection === true;
}

export const registerResponseDto = z.object({
  user: userDto,
  tenant: tenantSummaryDto,
});

export type RegisterResponseDto = z.infer<typeof registerResponseDto>;

export const sessionDto = z.object({
  id: z.uuid(),
  device_name: z.string(),
  ip_address: z.string(),
  user_agent: z.string(),
  created_at: z.iso.datetime({ offset: true }),
  expires_at: z.iso.datetime({ offset: true }),
  current: z.boolean(),
});

export type SessionDto = z.infer<typeof sessionDto>;

export const sessionListDto = z.array(sessionDto);
