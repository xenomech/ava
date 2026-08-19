import { queryOptions } from "@tanstack/react-query";

import { getFlow } from "./api";

export const flowQueries = {
  key: (flowType: string) => ["flow", flowType] as const,

  state: (flowType: string) =>
    queryOptions({
      queryKey: flowQueries.key(flowType),
      queryFn: ({ signal }) => getFlow(flowType, signal),
      staleTime: Infinity,
    }),
};
