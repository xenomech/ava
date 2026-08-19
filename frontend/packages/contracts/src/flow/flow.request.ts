import { z } from "zod";

export const submitStepRequest = z.object({
  data: z.unknown(),
});

export type SubmitStepRequest = z.infer<typeof submitStepRequest>;
