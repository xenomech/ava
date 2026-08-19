import type { QueryClient } from "@tanstack/react-query";
import { redirect, type ParsedLocation } from "@tanstack/react-router";

import { sessionQuery } from "./queries";

type GuardArgs = {
  context: { queryClient: QueryClient };
  location: ParsedLocation;
};

export async function requireAuth({ context, location }: GuardArgs): Promise<void> {
  const session = await context.queryClient.ensureQueryData(sessionQuery);

  if (!session) {
    throw redirect({ to: "/auth/login", search: { redirect: location.href } });
  }
}

export async function requireGuest({ context }: GuardArgs): Promise<void> {
  const session = await context.queryClient.ensureQueryData(sessionQuery);

  if (session) {
    throw redirect({ to: "/" });
  }
}
