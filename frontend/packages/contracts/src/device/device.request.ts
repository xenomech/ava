import { z } from "zod";

export const updateDeviceRequest = z.object({
  name: z.string().min(1).max(100).optional(),
  room: z.string().max(80).optional(),
});

export type UpdateDeviceRequest = z.infer<typeof updateDeviceRequest>;

export const DEVICE_ACTIONS = ["power", "brightness", "color_temp"] as const;
export const deviceActionSchema = z.enum(DEVICE_ACTIONS);
export type DeviceAction = z.infer<typeof deviceActionSchema>;

export const sendCommandRequest = z.object({
  action: deviceActionSchema,
  value: z.union([z.boolean(), z.number()]),
});

export type SendCommandRequest = z.infer<typeof sendCommandRequest>;
