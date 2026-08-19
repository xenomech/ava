import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { HeadContent, Outlet, createRootRouteWithContext } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";

import { AppShell } from "@/app/shell";
import type { RouterContext } from "@/app/router";
import { NotFound } from "@/shared/layouts/not-found";
import { RouteErrorBoundary } from "@/shared/layouts/route-error-boundary";

import "../index.css";

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
  notFoundComponent: NotFound,
  errorComponent: RouteErrorBoundary,
  head: () => ({
    meta: [
      { title: "Ava" },
      { name: "description", content: "Ava - multi-tenant hub control" },
    ],
    links: [{ rel: "icon", href: "/favicon.ico" }],
  }),
});

function RootComponent() {
  return (
    <>
      <HeadContent />
      <AppShell>
        <Outlet />
      </AppShell>
      <TanStackRouterDevtools position="bottom-left" />
      <ReactQueryDevtools buttonPosition="bottom-right" />
    </>
  );
}
