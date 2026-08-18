import { z } from "zod";

export const userDto = z.object({
  id: z.uuid(),
  email: z.email(),
  username: z.string(),
  name: z.string(),
  phone: z.string().optional(),
  email_verified: z.boolean(),
  email_verified_at: z.iso.datetime({ offset: true }).optional(),
  created_at: z.iso.datetime({ offset: true }),
});

export type UserDto = z.infer<typeof userDto>;
