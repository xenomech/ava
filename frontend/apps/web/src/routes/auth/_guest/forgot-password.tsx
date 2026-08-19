import { createFileRoute } from "@tanstack/react-router";

import { ForgotPasswordPage } from "@/modules/auth/pages/forgot-password-page";

export const Route = createFileRoute("/auth/_guest/forgot-password")({
  component: ForgotPasswordPage,
});
