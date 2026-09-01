import { queryOptions } from "@tanstack/react-query";

import { listTokens } from "./api";

export const tokenQueries = {
  all: () => ["tokens"] as const,

  list: () =>
    queryOptions({
      queryKey: [...tokenQueries.all(), "list"],
      queryFn: ({ signal }) => listTokens(signal),
    }),
};
