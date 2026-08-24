import { createFileRoute } from "@tanstack/react-router";

import { ActivatePage } from "@/modules/hub/pages/activate-page";

export const Route = createFileRoute("/_protected/settings/hubs")({
  component: ActivatePage,
});
