import { z } from "zod";
import { traitValueSchema } from "./device.dto";

export const updateDeviceRequest = z.object({
  name: z.string().min(1).max(100).optional(),
  room_id: z.uuid().optional(),
  clear_room: z.boolean().optional(),
  appliance: z.string().max(40).optional(),
});

export type UpdateDeviceRequest = z.infer<typeof updateDeviceRequest>;

export const sendCommandRequest = z.object({
  trait: z.string().min(1).max(64),
  value: traitValueSchema,
});

export type SendCommandRequest = z.infer<typeof sendCommandRequest>;

export const applyTargetRequest = z.object({
  device_id: z.uuid(),
  trait: z.string().min(1).max(64),
  value: traitValueSchema,
});

export type ApplyTargetRequest = z.infer<typeof applyTargetRequest>;

export const applyRequest = z.object({
  targets: z.array(applyTargetRequest).min(1).max(100),
});

export type ApplyRequest = z.infer<typeof applyRequest>;

export const applyResponse = z.object({
  applied: z.array(z.uuid()).default([]),
  skipped: z.array(z.object({ device_id: z.uuid(), reason: z.string() })).default([]),
});

export type ApplyResponse = z.infer<typeof applyResponse>;
