import { createFileRoute } from "@tanstack/react-router";

import { AcceptInvitePage } from "@/modules/auth/pages/accept-invite-page";

export const Route = createFileRoute("/auth/invite")({
  validateSearch: (search: Record<string, unknown>) => ({
    token: typeof search.token === "string" ? search.token : undefined,
  }),
  component: AcceptInvitePage,
});
