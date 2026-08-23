import { queryOptions } from "@tanstack/react-query";

import { listRooms } from "./api";

export const roomQueries = {
  all: () => ["rooms"] as const,
  list: () =>
    queryOptions({
      queryKey: [...roomQueries.all(), "list"] as const,
      queryFn: ({ signal }) => listRooms(signal),
    }),
};
