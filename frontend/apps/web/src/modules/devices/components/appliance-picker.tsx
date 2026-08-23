import { APPLIANCE_CHOICES, type DeviceDto } from "@ava/contracts";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@ava/ui";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { updateDevice } from "../api";
import { deviceQueries } from "../queries";

const NOTHING = "plug";

export function AppliancePicker({ device }: { device: DeviceDto }) {
  const queryClient = useQueryClient();

  const save = useMutation({
    mutationFn: (appliance: string) =>
      updateDevice(device.id, { appliance: appliance === NOTHING ? "" : appliance }),
    onSuccess: (updated) => {
      queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) =>
        current?.map((entry) => (entry.id === updated.id ? updated : entry)),
      );
    },
    onError: (error) =>
      toast.error(isApiError(error) ? error.message : "Could not save what is plugged in"),
  });

  return (
    <div className="grid gap-2.5">
      <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
        Plugged in
      </span>
      <Select
        value={device.appliance || NOTHING}
        onValueChange={(value) => save.mutate(value)}
        disabled={save.isPending}
      >
        <SelectTrigger aria-label="What is plugged in">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value={NOTHING}>Nothing specific</SelectItem>
          {APPLIANCE_CHOICES.map((choice) => (
            <SelectItem key={choice.value} value={choice.value}>
              {choice.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
