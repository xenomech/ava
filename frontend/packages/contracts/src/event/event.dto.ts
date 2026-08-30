import { z } from "zod";

import { deviceDto, type TraitValue } from "../device";

export const deviceStateEvent = z.object({
  type: z.literal("device.state"),
  device: deviceDto,
});

export const deviceListEvent = z.object({
  type: z.literal("device.list"),
  hub_id: z.uuid(),
  devices: z.array(deviceDto),
});

export const hubPresenceEvent = z.object({
  type: z.literal("hub.presence"),
  hub_id: z.uuid(),
  online: z.boolean(),
});

export const commandRejectedEvent = z.object({
  type: z.literal("command.rejected"),
  device_id: z.uuid(),
  reason: z.string(),
  message: z.string(),
});

export const avaEvent = z.discriminatedUnion("type", [
  deviceStateEvent,
  deviceListEvent,
  hubPresenceEvent,
  commandRejectedEvent,
]);

export type DeviceCommandFrame = {
  type: "device.command";
  device_id: string;
  trait: string;
  value: TraitValue;
};

export type CommandRejectedEvent = z.infer<typeof commandRejectedEvent>;
export type DeviceStateEvent = z.infer<typeof deviceStateEvent>;
export type DeviceListEvent = z.infer<typeof deviceListEvent>;
export type HubPresenceEvent = z.infer<typeof hubPresenceEvent>;
export type AvaEvent = z.infer<typeof avaEvent>;
