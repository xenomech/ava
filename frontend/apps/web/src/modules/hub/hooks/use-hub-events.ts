import type { AvaEvent, HubDto } from "@ava/contracts";
import { useQueryClient, type QueryClient } from "@tanstack/react-query";
import { useCallback } from "react";

import { useAvaEvent } from "@/shared/realtime";
import { hubQueries } from "../queries";

function setHubPresence(queryClient: QueryClient, hubID: string, online: boolean) {
  queryClient.setQueryData<HubDto[]>(hubQueries.list().queryKey, (current) =>
    current?.map((hub) => (hub.id === hubID && hub.online !== online ? { ...hub, online } : hub)),
  );
}

/** Keeps the hub cache in step with which hubs are actually answering. */
export function useHubEvents() {
  const queryClient = useQueryClient();

  useAvaEvent(
    useCallback(
      (event: AvaEvent) => {
        if (event.type === "hub.presence") setHubPresence(queryClient, event.hub_id, event.online);
      },
      [queryClient],
    ),
  );
}
