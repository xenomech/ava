import { avaEvent, type DeviceDto, type HubDto } from "@ava/contracts";
import { useQueryClient, type QueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import { toast } from "sonner";

import { deviceQueries, reconcile, reconcileAll, release } from "@/modules/devices";
import { hubQueries } from "@/modules/hub";
import { useAvaSocket } from "@/shared/realtime";

function byCreation(one: DeviceDto, other: DeviceDto) {
  return one.created_at.localeCompare(other.created_at);
}

/* Reconciled, not taken at face value. Anything this browser is still writing
   stays as the person left it until the house agrees or the hold expires. */
function replaceDevice(queryClient: QueryClient, device: DeviceDto) {
  const settled = reconcile(device);

  queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) =>
    current?.map((existing) => (existing.id === settled.id ? settled : existing)),
  );
}

function replaceHubDevices(queryClient: QueryClient, hubID: string, devices: DeviceDto[]) {
  queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) => {
    if (!current) return current;

    const untouched = current.filter((device) => device.hub_id !== hubID);
    const merged = [...untouched, ...reconcileAll(devices)].sort(byCreation);

    return JSON.stringify(merged) === JSON.stringify(current) ? current : merged;
  });
}

function setHubPresence(queryClient: QueryClient, hubID: string, online: boolean) {
  queryClient.setQueryData<HubDto[]>(hubQueries.list().queryKey, (current) =>
    current?.map((hub) => (hub.id === hubID && hub.online !== online ? { ...hub, online } : hub)),
  );
}

function apply(queryClient: QueryClient, raw: string) {
  const parsed = avaEvent.safeParse(JSON.parse(raw));
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
    case "command.rejected":
      /* The write is not ours to hold any more — let the refetch below say what
         the device is really doing. */
      release(event.device_id);
      toast.error(event.message);
      void queryClient.invalidateQueries({ queryKey: deviceQueries.all() });
      break;
  }
}

export function AvaEvents() {
  const queryClient = useQueryClient();
  const socket = useAvaSocket();

  useEffect(() => socket.subscribe((raw) => apply(queryClient, raw)), [socket, queryClient]);

  return null;
}
