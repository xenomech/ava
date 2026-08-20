import { z } from "zod";

export const DEVICE_STATUSES = ["online", "offline"] as const;
export const deviceStatusSchema = z.enum(DEVICE_STATUSES);
export type DeviceStatus = z.infer<typeof deviceStatusSchema>;

export const DEVICE_CAPABILITIES = ["brightness", "color_temp", "color"] as const;
export const deviceCapabilitySchema = z.enum(DEVICE_CAPABILITIES);
export type DeviceCapability = z.infer<typeof deviceCapabilitySchema>;

export const deviceLimitsDto = z.object({
  brightness_min: z.number().default(0),
  brightness_max: z.number().default(100),
  kelvin_min: z.number().optional(),
  kelvin_max: z.number().optional(),
});

export type DeviceLimitsDto = z.infer<typeof deviceLimitsDto>;

export const deviceStateDto = z.looseObject({
  power: z.boolean().default(false),
  brightness: z.number().optional(),
  color_temp: z.number().optional(),
  capabilities: z.array(deviceCapabilitySchema).default([]),
  limits: deviceLimitsDto.optional(),
  model: z.string().optional(),
  vendor: z.string().optional(),
  ip: z.string().optional(),
});

export type DeviceStateDto = z.infer<typeof deviceStateDto>;

export const deviceDto = z.object({
  id: z.uuid(),
  hub_id: z.uuid(),
  external_id: z.string(),
  name: z.string(),
  room: z.string(),
  kind: z.string(),
  status: deviceStatusSchema,
  last_seen_at: z.iso.datetime({ offset: true }).optional(),
  state: deviceStateDto,
  created_at: z.iso.datetime({ offset: true }),
});

export type DeviceDto = z.infer<typeof deviceDto>;
