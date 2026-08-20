import type { DeviceAction, DeviceDto } from "@ava/contracts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { sendCommand } from "./api";
import { deviceQueries } from "./queries";

const SETTLE_MS = 2_000;

export function useDevices() {
  const query = useQuery(deviceQueries.list());

  return {
    devices: query.data ?? [],
    isPending: query.isPending,
    isError: query.isError,
    error: query.error,
  };
}

type CommandInput = {
  device: DeviceDto;
  action: DeviceAction;
  value: boolean | number;
};

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

export function useDeviceCommand() {
  const queryClient = useQueryClient();
  const key = deviceQueries.list().queryKey;

  return useMutation({
    mutationFn: ({ device, action, value }: CommandInput) =>
      sendCommand(device.id, { action, value }),

    onMutate: async ({ device, action, value }: CommandInput) => {
      await queryClient.cancelQueries({ queryKey: deviceQueries.all() });

      const previous = queryClient.getQueryData<DeviceDto[]>(key);

      queryClient.setQueryData<DeviceDto[]>(key, (current) =>
        current?.map((entry) =>
          entry.id === device.id ? applyLocally(entry, action, value) : entry,
        ),
      );

      return { previous };
    },

    onError: (error, _input, context) => {
      if (context?.previous) queryClient.setQueryData(key, context.previous);

      toast.error(isApiError(error) ? error.message : "The hub did not accept that");
    },

    onSettled: () => {
      setTimeout(() => {
        void queryClient.invalidateQueries({ queryKey: deviceQueries.all() });
      }, SETTLE_MS);
    },
  });
}
