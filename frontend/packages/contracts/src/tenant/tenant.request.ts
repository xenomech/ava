import { z } from "zod";

import { tenantRoleSchema } from "./tenant.dto";

export const tenantSlugSchema = z
  .string()
  .min(3)
  .max(40)
  .regex(/^[a-z0-9]+(?:-[a-z0-9]+)*$/, "Use lowercase letters, numbers and hyphens");

export const createTenantRequest = z.object({
  name: z.string().min(1),
  slug: tenantSlugSchema,
});

export type CreateTenantRequest = z.infer<typeof createTenantRequest>;

export const updateTenantRequest = z.object({
  name: z.string().min(1),
});

export type UpdateTenantRequest = z.infer<typeof updateTenantRequest>;

export const inviteMemberRequest = z.object({
  email: z.email(),
  role: tenantRoleSchema,
});

export type InviteMemberRequest = z.infer<typeof inviteMemberRequest>;

export const updateMemberRoleRequest = z.object({
  role: tenantRoleSchema,
});

export type UpdateMemberRoleRequest = z.infer<typeof updateMemberRoleRequest>;
