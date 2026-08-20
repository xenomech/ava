import { createFileRoute } from "@tanstack/react-router";

import { ConsolePage } from "@/modules/devices/pages/console-page";

export const Route = createFileRoute("/_protected/")({
  validateSearch: (search: Record<string, unknown>): { device?: string } =>
    typeof search.device === "string" ? { device: search.device } : {},
  component: ConsolePage,
});
