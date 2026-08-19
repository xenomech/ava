import { createFileRoute } from "@tanstack/react-router";

import { RegisterPage } from "@/modules/auth/pages/register-page";

export const Route = createFileRoute("/auth/_guest/register")({
  component: RegisterPage,
});
