import { z } from "zod";

// The permissions a token may carry. Mirrors model.AllScopes on the server.
export const TOKEN_SCOPES = [
  "devices:read",
  "devices:write",
  "rooms:read",
  "rooms:write",
  "scenes:read",
  "scenes:write",
  "hubs:read",
  "hubs:write",
] as const;

export const tokenScopeSchema = z.enum(TOKEN_SCOPES);
export type TokenScope = z.infer<typeof tokenScopeSchema>;

export const apiTokenDto = z.object({
  id: z.uuid(),
  name: z.string(),
  scopes: z.array(tokenScopeSchema),
  last_used_at: z.string().nullable(),
  expires_at: z.string().nullable(),
  revoked_at: z.string().nullable(),
  created_at: z.string(),
});

export type ApiTokenDto = z.infer<typeof apiTokenDto>;

// The plaintext arrives once, on creation, and is never retrievable again.
export const createdApiTokenDto = z.object({
  token: apiTokenDto,
  value: z.string(),
});

export type CreatedApiTokenDto = z.infer<typeof createdApiTokenDto>;

export const tokenScopesDto = z.object({ scopes: z.array(tokenScopeSchema) });

export type TokenScopesDto = z.infer<typeof tokenScopesDto>;
