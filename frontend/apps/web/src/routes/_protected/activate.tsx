import { createFileRoute } from "@tanstack/react-router";

import { ActivatePage } from "@/modules/hub/pages/activate-page";

export const Route = createFileRoute("/_protected/activate")({
  validateSearch: (search: Record<string, unknown>): { code?: string } =>
    typeof search.code === "string" ? { code: search.code } : {},
  component: ActivatePage,
});
