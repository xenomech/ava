import { createFileRoute } from "@tanstack/react-router";

import { DesignPage } from "@/modules/design/pages/design-page";

export const Route = createFileRoute("/design")({
  component: DesignPage,
});
