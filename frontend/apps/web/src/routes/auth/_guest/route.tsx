import { createFileRoute, Outlet } from "@tanstack/react-router";

import { requireGuest } from "@/modules/auth";

export const Route = createFileRoute("/auth/_guest")({
  beforeLoad: requireGuest,
  component: Outlet,
});
