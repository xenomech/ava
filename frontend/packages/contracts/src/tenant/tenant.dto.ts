import { z } from "zod";

export const TENANT_ROLES = ["owner", "admin", "member"] as const;
export const tenantRoleSchema = z.enum(TENANT_ROLES);
export type TenantRole = z.infer<typeof tenantRoleSchema>;

export const MEMBERSHIP_STATUSES = ["invited", "active"] as const;
export const membershipStatusSchema = z.enum(MEMBERSHIP_STATUSES);
export type MembershipStatus = z.infer<typeof membershipStatusSchema>;

export const tenantSummaryDto = z.object({
  id: z.uuid(),
  name: z.string(),
  slug: z.string(),
  role: tenantRoleSchema,
});

export type TenantSummaryDto = z.infer<typeof tenantSummaryDto>;

export const tenantDto = z.object({
  id: z.uuid(),
  name: z.string(),
  slug: z.string(),
  created_at: z.iso.datetime({ offset: true }),
});

export type TenantDto = z.infer<typeof tenantDto>;

export const memberDto = z.object({
  user_id: z.uuid(),
  email: z.email(),
  name: z.string(),
  role: tenantRoleSchema,
  status: membershipStatusSchema,
  joined_at: z.iso.datetime({ offset: true }).optional(),
});

export type MemberDto = z.infer<typeof memberDto>;
