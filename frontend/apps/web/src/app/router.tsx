import type { QueryClient } from "@tanstack/react-query";
import { createRouter } from "@tanstack/react-router";

import Loader from "@/shared/components/loader";
import { routeTree } from "../routeTree.gen";

export type RouterContext = {
  queryClient: QueryClient;
};

export function buildRouter(queryClient: QueryClient) {
  return createRouter({
    routeTree,
    defaultPreload: "intent",
    scrollRestoration: true,
    defaultPendingComponent: () => <Loader />,
    context: { queryClient },
  });
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof buildRouter>;
  }
}
