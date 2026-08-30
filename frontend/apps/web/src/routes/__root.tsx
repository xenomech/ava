import { HeadContent, Outlet, createRootRouteWithContext } from "@tanstack/react-router";
import { Suspense, lazy } from "react";

import type { RouterContext } from "@/app/router";
import { UpdatePrompt } from "@/app/update-prompt";
import { NotFound } from "@/shared/layouts/not-found";
import { RouteErrorBoundary } from "@/shared/layouts/route-error-boundary";

import "../index.css";

/* Dev-only, and lazy even then: the root route is in every page's first chunk,
   and a static import here would ship both devtool bundles to production. */
const RouterDevtools = import.meta.env.DEV
  ? lazy(() =>
      import("@tanstack/react-router-devtools").then((module) => ({
        default: module.TanStackRouterDevtools,
      })),
    )
  : () => null;

const QueryDevtools = import.meta.env.DEV
  ? lazy(() =>
      import("@tanstack/react-query-devtools").then((module) => ({
        default: module.ReactQueryDevtools,
      })),
    )
  : () => null;

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
  notFoundComponent: NotFound,
  errorComponent: RouteErrorBoundary,
  head: () => ({
    meta: [{ title: "Ava" }, { name: "description", content: "Ava - multi-tenant hub control" }],
    links: [{ rel: "icon", href: "/favicon.ico" }],
  }),
});

function RootComponent() {
  return (
    <>
      <HeadContent />
      <UpdatePrompt />
      <Outlet />
      <Suspense fallback={null}>
        <RouterDevtools position="bottom-left" />
        <QueryDevtools buttonPosition="bottom-right" />
      </Suspense>
    </>
  );
}
