import { queryOptions } from "@tanstack/react-query";

import { getCurrentTenant, listMembers, listMyTenants } from "./api";

export const tenantQueries = {
  all: () => ["tenant"] as const,

  mine: () =>
    queryOptions({
      queryKey: [...tenantQueries.all(), "mine"],
      queryFn: ({ signal }) => listMyTenants(signal),
    }),

  current: () =>
    queryOptions({
      queryKey: [...tenantQueries.all(), "current"],
      queryFn: ({ signal }) => getCurrentTenant(signal),
    }),

  members: () =>
    queryOptions({
      queryKey: [...tenantQueries.all(), "members"],
      queryFn: ({ signal }) => listMembers(signal),
    }),
};
