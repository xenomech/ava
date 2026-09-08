import { z } from "zod";

import { tokenScopeSchema } from "./token.dto";

export const createApiTokenRequest = z.object({
  name: z.string().min(1, "Name is required").max(80),
  scopes: z.array(tokenScopeSchema).min(1, "Pick at least one permission"),
  // Omitted for a token that never expires.
  expires_in_days: z.number().int().min(1).max(3650).optional(),
});

export type CreateApiTokenRequest = z.infer<typeof createApiTokenRequest>;
