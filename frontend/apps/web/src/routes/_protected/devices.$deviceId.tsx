import { createFileRoute } from "@tanstack/react-router";

import { DevicePage } from "@/modules/devices/pages/device-page";

export const Route = createFileRoute("/_protected/devices/$deviceId")({
  component: DevicePage,
});
