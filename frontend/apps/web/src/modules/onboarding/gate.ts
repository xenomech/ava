import type { QueryClient } from "@tanstack/react-query";

import { requireFlowCompleted } from "@/shared/flow";
import { ONBOARDING_FLOW } from "./flow";

export function requireOnboarded({ queryClient }: { queryClient: QueryClient }): Promise<void> {
  return requireFlowCompleted({
    queryClient,
    flowType: ONBOARDING_FLOW,
    redirectTo: "/onboarding",
  });
}
