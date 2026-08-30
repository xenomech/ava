import { z } from "zod";

export const createRoomRequest = z.object({
  name: z.string().min(1).max(80),
});

export type CreateRoomRequest = z.infer<typeof createRoomRequest>;

export const updateRoomRequest = z.object({
  name: z.string().min(1).max(80).optional(),
  position: z.number().min(0).optional(),
});

export type UpdateRoomRequest = z.infer<typeof updateRoomRequest>;
