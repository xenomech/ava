import { createFileRoute } from "@tanstack/react-router";

import { requireGuest } from "@/modules/auth";
import { LoginPage } from "@/modules/auth/pages/login-page";

export const Route = createFileRoute("/auth/login")({
  beforeLoad: requireGuest,
  component: LoginPage,
});
