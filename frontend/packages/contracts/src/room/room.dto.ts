import { z } from "zod";

export const roomDto = z.object({
  id: z.uuid(),
  name: z.string(),
  position: z.number(),
});

export type RoomDto = z.infer<typeof roomDto>;
