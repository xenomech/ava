import type { AvaEvent, DeviceDto } from "@ava/contracts";
import { useQueryClient, type QueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import { toast } from "sonner";

import { useAvaEvent } from "@/shared/realtime";
import { reconcile, reconcileAll, release } from "../lib/optimistic";
import { deviceQueries } from "../queries";

function byCreation(one: DeviceDto, other: DeviceDto) {
  return one.created_at.localeCompare(other.created_at);
}

// Reconciled, so anything this browser is still writing stays as the person left it.
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

function apply(queryClient: QueryClient, event: AvaEvent) {
  switch (event.type) {
    case "device.state":
      replaceDevice(queryClient, event.device);
      break;
    case "device.list":
      replaceHubDevices(queryClient, event.hub_id, event.devices);
      break;
    case "command.rejected":
      // The write is not ours to hold any more; the refetch below says what is really true.
      release(event.device_id);
      toast.error(event.message);
      void queryClient.invalidateQueries({ queryKey: deviceQueries.all() });
      break;
  }
}

/** Keeps the device cache in step with what the house reports. */
export function useDeviceEvents() {
  const queryClient = useQueryClient();

  useAvaEvent(useCallback((event: AvaEvent) => apply(queryClient, event), [queryClient]));
}
