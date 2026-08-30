import { createFileRoute, redirect } from "@tanstack/react-router";

import { LooseDevicePage } from "@/modules/devices/pages/device-page";
import { deviceQueries } from "@/modules/devices/queries";

/** A doorway: a device in a room opens that room, and only a roomless device renders here. */
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
