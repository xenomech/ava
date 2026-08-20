import { z } from "zod";

export const HUB_STATUSES = ["active", "revoked"] as const;
export const hubStatusSchema = z.enum(HUB_STATUSES);
export type HubStatus = z.infer<typeof hubStatusSchema>;

export const hubDto = z.object({
  id: z.uuid(),
  name: z.string(),
  status: hubStatusSchema,
  online: z.boolean().default(false),
  last_seen_at: z.iso.datetime({ offset: true }).optional(),
  created_at: z.iso.datetime({ offset: true }),
});

export type HubDto = z.infer<typeof hubDto>;
