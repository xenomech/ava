import { createFileRoute } from "@tanstack/react-router";

import { requireGuest } from "@/modules/auth";
import { RegisterPage } from "@/modules/auth/pages/register-page";

export const Route = createFileRoute("/auth/register")({
  beforeLoad: requireGuest,
  component: RegisterPage,
});
