import { redirect, type LinkProps } from "@tanstack/react-router";
import type { QueryClient } from "@tanstack/react-query";

import { isApiError } from "@/config/http/request";
import { flowQueries } from "./queries";

export async function requireFlowCompleted({
  queryClient,
  flowType,
  redirectTo,
}: {
  queryClient: QueryClient;
  flowType: string;
  redirectTo: LinkProps["to"];
}): Promise<void> {
  let flow;

  try {
    flow = await queryClient.ensureQueryData(flowQueries.state(flowType));
  } catch (error: unknown) {
    /* No flow on record means nothing to complete. Anything else — a network
       failure, a 500 — must not silently pass the gate as if it had been
       checked; it rethrows for the route's error boundary. */
    if (isApiError(error) && error.status === 404) return;

    throw error;
  }

  if (flow.status === "completed") return;

  throw redirect({ to: redirectTo, replace: true });
}
