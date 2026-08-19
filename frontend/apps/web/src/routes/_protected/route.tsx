import { createFileRoute } from "@tanstack/react-router";

import { AppShell } from "@/app/shell";
import { requireAuth } from "@/modules/auth";
import { requireOnboarded } from "@/modules/onboarding";

export const Route = createFileRoute("/_protected")({
  beforeLoad: async ({ context, location }) => {
    await requireAuth({ context, location });
    await requireOnboarded(context);
  },
  component: AppShell,
});
