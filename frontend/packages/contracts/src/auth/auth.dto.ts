import { z } from "zod";

import { tenantSummaryDto } from "../tenant/tenant.dto";
import { userDto } from "../user/user.dto";

export const tokensDto = z.object({
  expires_in: z.number().int(),
  access_token: z.string().optional(),
  refresh_token: z.string().optional(),
});

export type TokensDto = z.infer<typeof tokensDto>;

export const authenticatedDto = z.object({
  user: userDto,
  tenant: tenantSummaryDto,
  tokens: tokensDto,
});

export const authResponseDto = authenticatedDto;

export type AuthenticatedDto = z.infer<typeof authenticatedDto>;
export type AuthResponseDto = z.infer<typeof authResponseDto>;

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
