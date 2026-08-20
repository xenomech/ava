import { queryOptions } from "@tanstack/react-query";

import { listHubs } from "./api";

export const hubQueries = {
  all: () => ["hub"] as const,

  list: () =>
    queryOptions({
      queryKey: [...hubQueries.all(), "list"],
      queryFn: ({ signal }) => listHubs(signal),
    }),
};
