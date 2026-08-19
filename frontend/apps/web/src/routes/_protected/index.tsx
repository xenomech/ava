import { createFileRoute } from "@tanstack/react-router";

import { ConsolePage } from "@/modules/devices/pages/console-page";

export const Route = createFileRoute("/_protected/")({
  component: ConsolePage,
});
