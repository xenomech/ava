import type { DeviceDto } from "@ava/contracts";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@ava/ui";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { useRooms } from "@/modules/rooms";
import { updateDevice } from "../api";
import { deviceQueries } from "../queries";

const UNASSIGNED = "none";

export function RoomPicker({ device, label = true }: { device: DeviceDto; label?: boolean }) {
  const queryClient = useQueryClient();
  const { rooms } = useRooms();

  const save = useMutation({
    mutationFn: (roomID: string) =>
      updateDevice(device.id, roomID === UNASSIGNED ? { clear_room: true } : { room_id: roomID }),
    onSuccess: (updated) => {
      queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) =>
        current?.map((entry) => (entry.id === updated.id ? updated : entry)),
      );
    },
    onError: (error) =>
      toast.error(isApiError(error) ? error.message : "Could not move the device"),
  });

  return (
    <div className="grid gap-2.5">
      {label ? (
        <span className="text-caption font-semibold uppercase tracking-caps text-subtle">Room</span>
      ) : null}
      <Select
        value={device.room_id ?? UNASSIGNED}
        onValueChange={(value) => save.mutate(value)}
        disabled={save.isPending}
      >
        <SelectTrigger aria-label={`Room for ${device.name}`}>
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value={UNASSIGNED}>No room</SelectItem>
          {rooms.map((room) => (
            <SelectItem key={room.id} value={room.id}>
              {room.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
