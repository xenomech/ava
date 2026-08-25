import {
  TRAIT_BRIGHTNESS,
  TRAIT_COLOR,
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  supports,
  type ApplyTargetRequest,
  type DeviceDto,
  type TraitValue,
} from "@ava/contracts";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import { toast } from "sonner";

import { useAvaSocket } from "@/shared/realtime";
import { isApiError } from "@/config/http/request";
import { applyTargets, sendCommand } from "./api";
import { deviceQueries } from "./queries";

export function useDevices() {
  const query = useQuery(deviceQueries.list());

  return {
    devices: query.data ?? [],
    isPending: query.isPending,
    isError: query.isError,
    error: query.error,
  };
}

function applyLocally(device: DeviceDto, trait: string, value: TraitValue): DeviceDto {
  const state = { ...device.state, [trait]: value };

  if (trait === TRAIT_BRIGHTNESS && typeof value === "number") {
    state[TRAIT_POWER] = value > 0;
  }

  /* A light shows a colour or a temperature, never both, and the hub clears
     whichever one you did not set. The optimistic copy has to say the same
     thing: while it claimed a light had both, the colour panel read the stale
     temperature, stayed on the White tab, and so refused to switch back to it. */
  if (trait === TRAIT_COLOR) delete state[TRAIT_COLOR_TEMP];
  if (trait === TRAIT_COLOR_TEMP) delete state[TRAIT_COLOR];

  return { ...device, state };
}

export function useDeviceControl() {
  const queryClient = useQueryClient();
  const socket = useAvaSocket();

  return useCallback(
    (device: DeviceDto, trait: string, value: TraitValue) => {
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
        toast.error(isApiError(error) ? error.message : "The hub did not accept that");
        void queryClient.invalidateQueries({ queryKey: deviceQueries.all() });
      });
    },
    [queryClient, socket],
  );
}

/**
 * Push a batch of trait writes, showing them locally before the hub answers.
 *
 * Shared by the room switch and by scenes, because a scene is nothing more than
 * a batch of writes someone saved earlier — there is no second code path for
 * replaying one, and so no second place for the optimistic copy to drift.
 */
export function useApplyTargets() {
  const queryClient = useQueryClient();

  return useCallback(
    async (targets: ApplyTargetRequest[]) => {
      if (targets.length === 0) {
        toast.error("Nothing here can be switched");

        return;
      }

      const byDevice = new Map<string, ApplyTargetRequest[]>();
      for (const target of targets) {
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
        const result = await applyTargets({ targets });

        if (result.skipped.length > 0) {
          toast.warning(`${result.applied.length} changed, ${result.skipped.length} skipped`);
        }
      } catch (error: unknown) {
        toast.error(isApiError(error) ? error.message : "The hub did not accept that");
        void queryClient.invalidateQueries({ queryKey: deviceQueries.all() });
      }
    },
    [queryClient],
  );
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
