import {
  TRAIT_BRIGHTNESS,
  TRAIT_POWER,
  supports,
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

export function useRoomPower() {
  const queryClient = useQueryClient();

  return useCallback(
    async (devices: DeviceDto[], on: boolean) => {
      const targets = devices
        .filter((device) => device.status !== "offline" && supports(device, TRAIT_POWER))
        .map((device) => ({ device_id: device.id, trait: TRAIT_POWER, value: on }));

      if (targets.length === 0) {
        toast.error("Nothing here can be switched");

        return;
      }

      const ids = new Set(targets.map((target) => target.device_id));

      queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) =>
        current?.map((entry) => (ids.has(entry.id) ? applyLocally(entry, TRAIT_POWER, on) : entry)),
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
