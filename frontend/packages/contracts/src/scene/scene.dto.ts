import { z } from "zod";
import { traitValueSchema } from "../device/device.dto";

export const sceneTargetDto = z.object({
  device_id: z.uuid(),
  trait: z.string(),
  value: traitValueSchema,
});

export type SceneTargetDto = z.infer<typeof sceneTargetDto>;

export const sceneDto = z.object({
  id: z.uuid(),
  room_id: z.uuid(),
  name: z.string(),
  position: z.number(),
  targets: z.array(sceneTargetDto).default([]),
});

export type SceneDto = z.infer<typeof sceneDto>;
