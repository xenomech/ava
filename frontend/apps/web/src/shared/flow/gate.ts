import { redirect, type LinkProps } from "@tanstack/react-router";
import type { QueryClient } from "@tanstack/react-query";

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
  } catch {
    return;
  }

  if (flow.status === "completed") return;

  throw redirect({ to: redirectTo, replace: true });
}
