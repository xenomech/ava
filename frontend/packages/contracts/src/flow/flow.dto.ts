import { z } from "zod";

import { tenantRoleSchema } from "../tenant/tenant.dto";

export const FLOW_STATUSES = ["pending", "in_progress", "completed", "failed"] as const;
export const flowStatusSchema = z.enum(FLOW_STATUSES);
export type FlowStatus = z.infer<typeof flowStatusSchema>;

export const FLOW_STEP_STATUSES = [
  "pending",
  "in_progress",
  "completed",
  "skipped",
  "failed",
] as const;
export const flowStepStatusSchema = z.enum(FLOW_STEP_STATUSES);
export type FlowStepStatus = z.infer<typeof flowStepStatusSchema>;

export const flowStepDto = z.object({
  id: z.string(),
  title: z.string(),
  description: z.string(),
  status: flowStepStatusSchema,
  skippable: z.boolean(),
  data: z.unknown(),
  errors: z.record(z.string(), z.string()),
});

export type FlowStepDto = z.infer<typeof flowStepDto>;

export const flowStateDto = z.object({
  flow_type: z.string(),
  status: flowStatusSchema,
  current_step: z.string(),
  steps: z.array(flowStepDto),
  metadata: z.unknown(),
});

export type FlowStateDto = z.infer<typeof flowStateDto>;

export const onboardingMetadataDto = z.object({
  invite_roles: z.array(z.object({ value: tenantRoleSchema, label: z.string() })),
});

export type OnboardingMetadataDto = z.infer<typeof onboardingMetadataDto>;
