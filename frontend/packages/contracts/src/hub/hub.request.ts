import { z } from "zod";

export const activateHubRequest = z.object({
  user_code: z.string().min(1, "Enter the code shown by your hub"),
});

export type ActivateHubRequest = z.infer<typeof activateHubRequest>;
