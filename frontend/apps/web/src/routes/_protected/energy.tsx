import { createFileRoute } from "@tanstack/react-router";

import { EnergyPage } from "@/modules/devices/pages/energy-page";

export const Route = createFileRoute("/_protected/energy")({
  component: EnergyPage,
});
