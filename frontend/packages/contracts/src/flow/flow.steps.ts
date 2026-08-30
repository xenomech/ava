import { z } from "zod";

export const homeStepData = z.object({
  name: z.string().min(1, "Give your home a name"),
});

// The hub step carries no fields. Pairing runs through the hub activation
// endpoint; submitting the step only asks the server to confirm a hub arrived.
export const hubStepData = z.object({});

export type HomeStepData = z.infer<typeof homeStepData>;
export type HubStepData = z.infer<typeof hubStepData>;

export const ONBOARDING_STEP_SCHEMAS = {
  home: homeStepData,
  hub: hubStepData,
} as const;

export type OnboardingStepId = keyof typeof ONBOARDING_STEP_SCHEMAS;
