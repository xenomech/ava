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

/**
 * A device's readings, with anything unreadable dropped rather than fatal.
 *
 * Traits arrive as null when a device retires one — a bulb given a colour stops
 * reporting a temperature — and a strict record turned that single null into a
 * validation failure for the entire request. Every device in the house came
 * back empty because one bulb had changed colour, which is a spectacular
 * penalty for a value nothing was going to read anyway.
 *
 * So a null is dropped and the device keeps its other traits. Everything else
 * is still held to the contract: this is leniency about absence, not about
 * shape.
 */
export const deviceStateDto = z
  .record(z.string(), z.union([traitValueSchema, z.null()]))
  .transform(
    (state) =>
      Object.fromEntries(Object.entries(state).filter(([, value]) => value !== null)) as Record<
        string,
        TraitValue
      >,
  );

export type DeviceStateDto = Record<string, TraitValue>;

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
