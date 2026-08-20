import { avaEvent, type DeviceDto, type HubDto } from "@ava/contracts";
import { env } from "@ava/env/web";
import { useQueryClient, type QueryClient } from "@tanstack/react-query";
import { useEffect } from "react";

import { deviceQueries } from "@/modules/devices";
import { hubQueries } from "@/modules/hub";

function byCreation(one: DeviceDto, other: DeviceDto) {
  return one.created_at.localeCompare(other.created_at);
}

function replaceDevice(queryClient: QueryClient, device: DeviceDto) {
  queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) =>
    current?.map((existing) => (existing.id === device.id ? device : existing)),
  );
}

function replaceHubDevices(queryClient: QueryClient, hubID: string, devices: DeviceDto[]) {
  queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) => {
    if (!current) return current;

    const untouched = current.filter((device) => device.hub_id !== hubID);
    const merged = [...untouched, ...devices].sort(byCreation);

    return JSON.stringify(merged) === JSON.stringify(current) ? current : merged;
  });
}

function setHubPresence(queryClient: QueryClient, hubID: string, online: boolean) {
  queryClient.setQueryData<HubDto[]>(hubQueries.list().queryKey, (current) =>
    current?.map((hub) => (hub.id === hubID && hub.online !== online ? { ...hub, online } : hub)),
  );
}

export function useAvaStream() {
  const queryClient = useQueryClient();

  useEffect(() => {
    const source = new EventSource(`${env.VITE_API_URL}/events`, { withCredentials: true });

    source.onmessage = (message) => {
      const parsed = avaEvent.safeParse(JSON.parse(message.data));
      if (!parsed.success) return;

      const event = parsed.data;

      switch (event.type) {
        case "device.state":
          replaceDevice(queryClient, event.device);
          break;
        case "device.list":
          replaceHubDevices(queryClient, event.hub_id, event.devices);
          break;
        case "hub.presence":
          setHubPresence(queryClient, event.hub_id, event.online);
          break;
      }
    };

    return () => source.close();
  }, [queryClient]);
}
