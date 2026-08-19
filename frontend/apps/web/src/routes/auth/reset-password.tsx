import { createFileRoute } from "@tanstack/react-router";

import { ResetPasswordPage } from "@/modules/auth/pages/reset-password-page";

export const Route = createFileRoute("/auth/reset-password")({
  validateSearch: (search: Record<string, unknown>) => ({
    token: typeof search.token === "string" ? search.token : undefined,
  }),
  component: ResetPasswordPage,
});
