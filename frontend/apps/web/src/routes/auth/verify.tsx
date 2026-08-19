import { createFileRoute } from "@tanstack/react-router";

import { VerifyEmailPage } from "@/modules/auth/pages/verify-email-page";

export const Route = createFileRoute("/auth/verify")({
  validateSearch: (search: Record<string, unknown>) => ({
    token: typeof search.token === "string" ? search.token : undefined,
  }),
  component: VerifyEmailPage,
});
