import { createFileRoute } from "@tanstack/react-router";

import { AppShell } from "@/app/shell";
import { requireAuth } from "@/modules/auth";

export const Route = createFileRoute("/_protected")({
  beforeLoad: requireAuth,
  component: AppShell,
});
