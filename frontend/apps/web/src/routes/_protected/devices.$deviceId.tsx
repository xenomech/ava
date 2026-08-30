import { createFileRoute, redirect } from "@tanstack/react-router";

import { LooseDevicePage } from "@/modules/devices/pages/device-page";
import { deviceQueries } from "@/modules/devices/queries";

/**
 * A device's own address.
 *
 * Devices are reached through their room, so this is mostly a doorway: if the
 * device belongs to a room, it opens that room with the device's controls
 * already out. Only a device with no room renders here, because there is no
 * room to open.
 */
export const Route = createFileRoute("/_protected/devices/$deviceId")({
  beforeLoad: async ({ context, params }) => {
    const devices = await context.queryClient.ensureQueryData(deviceQueries.list());
    const device = devices.find((entry) => entry.id === params.deviceId);

    if (device?.room_id) {
      throw redirect({
        to: "/rooms/$roomId",
        params: { roomId: device.room_id },
        search: { device: device.id },
        replace: true,
      });
    }
  },
  component: LooseDevicePage,
});
