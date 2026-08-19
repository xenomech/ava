import { createFileRoute } from "@tanstack/react-router";

import { SettingsPage } from "@/modules/tenant/pages/settings-page";

export const Route = createFileRoute("/_protected/settings/")({
  component: SettingsPage,
});
