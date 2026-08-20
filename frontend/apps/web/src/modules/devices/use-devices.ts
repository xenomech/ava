import type { DeviceAction, DeviceDto } from "@ava/contracts";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import { toast } from "sonner";

import { useAvaSocket } from "@/shared/realtime";
import { isApiError } from "@/config/http/request";
import { sendCommand } from "./api";
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

function applyLocally(device: DeviceDto, action: DeviceAction, value: boolean | number): DeviceDto {
  const state = { ...device.state };

  if (action === "power" && typeof value === "boolean") state.power = value;
  if (action === "brightness" && typeof value === "number") {
    state.brightness = value;
    state.power = value > 0;
  }
  if (action === "color_temp" && typeof value === "number") state.color_temp = value;

  return { ...device, state };
}

export function useDeviceControl() {
  const queryClient = useQueryClient();
  const socket = useAvaSocket();

  return useCallback(
    (device: DeviceDto, action: DeviceAction, value: boolean | number) => {
      queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) =>
        current?.map((entry) =>
          entry.id === device.id ? applyLocally(entry, action, value) : entry,
        ),
      );

      const sent = socket.send({
        type: "device.command",
        device_id: device.id,
        action,
        value,
      });

      if (sent) return;

      void sendCommand(device.id, { action, value }).catch((error: unknown) => {
        toast.error(isApiError(error) ? error.message : "The hub did not accept that");
        void queryClient.invalidateQueries({ queryKey: deviceQueries.all() });
      });
    },
    [queryClient, socket],
  );
}
