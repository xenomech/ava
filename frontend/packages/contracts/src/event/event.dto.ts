import { z } from "zod";

import { deviceDto } from "../device";

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

export const avaEvent = z.discriminatedUnion("type", [
  deviceStateEvent,
  deviceListEvent,
  hubPresenceEvent,
]);

export type DeviceStateEvent = z.infer<typeof deviceStateEvent>;
export type DeviceListEvent = z.infer<typeof deviceListEvent>;
export type HubPresenceEvent = z.infer<typeof hubPresenceEvent>;
export type AvaEvent = z.infer<typeof avaEvent>;
