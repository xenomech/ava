import { z } from "zod";

import { tenantRoleSchema } from "../tenant/tenant.dto";

export const profileStepData = z.object({
  name: z.string().min(1, "Name is required"),
  phone: z.string().optional(),
});

export const workspaceStepData = z.object({
  name: z.string().min(1, "Workspace name is required"),
});

export const inviteTeamStepData = z.object({
  emails: z.array(z.email()),
  role: tenantRoleSchema.optional(),
});

export type ProfileStepData = z.infer<typeof profileStepData>;
export type WorkspaceStepData = z.infer<typeof workspaceStepData>;
export type InviteTeamStepData = z.infer<typeof inviteTeamStepData>;

export const ONBOARDING_STEP_SCHEMAS = {
  profile: profileStepData,
  workspace: workspaceStepData,
  invite_team: inviteTeamStepData,
} as const;

export type OnboardingStepId = keyof typeof ONBOARDING_STEP_SCHEMAS;
