import {
  TRAIT_POWER,
  supports,
  type ApplyResponse,
  type ApplyTargetRequest,
  type DeviceDto,
  type TraitValue,
} from "@ava/contracts";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import { toast } from "sonner";

import { useAvaSocket } from "@/shared/realtime";
import { isApiError } from "@/config/http/request";
import { applyTargets, sendCommand } from "../api";
import { applyLocally, claim, release } from "../lib/optimistic";
import { deviceQueries } from "../queries";

// One shared empty list, so a pending query does not hand consumers a fresh [] each render.
const NO_DEVICES: DeviceDto[] = [];

export function useDevices() {
  const query = useQuery(deviceQueries.list());

  return {
    devices: query.data ?? NO_DEVICES,
    isPending: query.isPending,
    isError: query.isError,
    error: query.error,
  };
}

export function useDeviceControl() {
  const queryClient = useQueryClient();
  const socket = useAvaSocket();

  return useCallback(
    (device: DeviceDto, trait: string, value: TraitValue) => {
      claim(device.id, trait, value);

      queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) =>
        current?.map((entry) =>
          entry.id === device.id ? applyLocally(entry, trait, value) : entry,
        ),
      );

      const sent = socket.send({
        type: "device.command",
        device_id: device.id,
        trait,
        value,
      });

      if (sent) return;

      void sendCommand(device.id, { trait, value }).catch((error: unknown) => {
        release(device.id, trait);
        toast.error(isApiError(error) ? error.message : "The hub did not accept that");
        void queryClient.invalidateQueries({ queryKey: deviceQueries.all() });
      });
    },
    [queryClient, socket],
  );
}

/** Push a batch of trait writes, shown locally before the hub answers. */
/**
 * Paint a batch of writes locally, then send them however the caller says.
 *
 * The optimistic half is the same whether the writes came from a slider or from a saved scene;
 * only the request differs, so the sender is passed in rather than assumed.
 */
function useOptimisticApply() {
  const queryClient = useQueryClient();

  return useCallback(
    async (targets: ApplyTargetRequest[], send: () => Promise<ApplyResponse>) => {
      if (targets.length === 0) {
        toast.error("Nothing here can be switched");

        return;
      }

      const byDevice = new Map<string, ApplyTargetRequest[]>();
      for (const target of targets) {
        claim(target.device_id, target.trait, target.value);
        byDevice.set(target.device_id, [...(byDevice.get(target.device_id) ?? []), target]);
      }

      queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) =>
        current?.map((entry) =>
          (byDevice.get(entry.id) ?? []).reduce(
            (device, target) => applyLocally(device, target.trait, target.value),
            entry,
          ),
        ),
      );

      try {
        const result = await send();

        if (result.skipped.length > 0) {
          toast.warning(`${result.applied.length} changed, ${result.skipped.length} skipped`);
        }
      } catch (error: unknown) {
        for (const target of targets) release(target.device_id, target.trait);

        toast.error(isApiError(error) ? error.message : "The hub did not accept that");
        void queryClient.invalidateQueries({ queryKey: deviceQueries.all() });
      }
    },
    [queryClient],
  );
}

/**
 * Push a batch of trait writes, showing them locally before the hub answers.
 *
 * Shared by the room switch and by anything else that changes several devices at once.
 */
export function useApplyTargets() {
  const optimistic = useOptimisticApply();

  return useCallback(
    (targets: ApplyTargetRequest[]) => optimistic(targets, () => applyTargets({ targets })),
    [optimistic],
  );
}

/** The same optimistic paint, for a write the server composes rather than the client. */
export function useOptimisticSend() {
  return useOptimisticApply();
}

export function useRoomPower() {
  const apply = useApplyTargets();

  return useCallback(
    (devices: DeviceDto[], on: boolean) =>
      apply(
        devices
          .filter((device) => device.status !== "offline" && supports(device, TRAIT_POWER))
          .map((device) => ({ device_id: device.id, trait: TRAIT_POWER, value: on })),
      ),
    [apply],
  );
}
