import { queryOptions } from "@tanstack/react-query";

import { listDevices } from "./api";

const SAFETY_NET_MS = 120_000;

export const deviceQueries = {
  all: () => ["device"] as const,

  list: () =>
    queryOptions({
      queryKey: [...deviceQueries.all(), "list"],
      queryFn: ({ signal }) => listDevices(signal),
      refetchInterval: SAFETY_NET_MS,
    }),
};
