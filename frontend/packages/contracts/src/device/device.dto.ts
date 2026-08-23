import { z } from "zod";

export const DEVICE_STATUSES = ["online", "offline"] as const;
export const deviceStatusSchema = z.enum(DEVICE_STATUSES);
export type DeviceStatus = z.infer<typeof deviceStatusSchema>;

export const TRAIT_KINDS = ["bool", "number", "enum", "color"] as const;
export const traitKindSchema = z.enum(TRAIT_KINDS);
export type TraitKind = z.infer<typeof traitKindSchema>;

export const TRAIT_ACCESS = ["r", "rw"] as const;
export const traitAccessSchema = z.enum(TRAIT_ACCESS);
export type TraitAccess = z.infer<typeof traitAccessSchema>;

export const TRAIT_POWER = "power";
export const TRAIT_BRIGHTNESS = "brightness";
export const TRAIT_COLOR_TEMP = "color_temp";
export const TRAIT_COLOR = "color";

export const capabilityDto = z.object({
  trait: z.string(),
  kind: traitKindSchema,
  access: traitAccessSchema,
  min: z.number().optional(),
  max: z.number().optional(),
  step: z.number().optional(),
  unit: z.string().optional(),
  values: z.array(z.string()).optional(),
});

export type CapabilityDto = z.infer<typeof capabilityDto>;

export const traitValueSchema = z.union([z.boolean(), z.number(), z.string()]);
export type TraitValue = z.infer<typeof traitValueSchema>;

export const deviceStateDto = z.record(z.string(), traitValueSchema);
export type DeviceStateDto = z.infer<typeof deviceStateDto>;

export const deviceDto = z.object({
  id: z.uuid(),
  hub_id: z.uuid(),
  external_id: z.string(),
  name: z.string(),
  room_id: z.uuid().optional(),
  room: z.string().default(""),
  appliance: z.string().default(""),
  kind: z.string(),
  vendor: z.string().optional(),
  model: z.string().optional(),
  parent: z.string().optional(),
  status: deviceStatusSchema,
  last_seen_at: z.iso.datetime({ offset: true }).optional(),
  capabilities: z.array(capabilityDto).default([]),
  state: deviceStateDto.default({}),
  created_at: z.iso.datetime({ offset: true }),
});

export type DeviceDto = z.infer<typeof deviceDto>;
