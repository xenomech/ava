import type { RoomDto } from "@ava/contracts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { createRoom, deleteRoom, updateRoom } from "./api";
import { roomQueries } from "./queries";

function reason(error: unknown, fallback: string) {
  return isApiError(error) ? error.message : fallback;
}

/** The list with one room shifted a step, for the reorder mutation. */
export function moved<T>(list: T[], index: number, direction: -1 | 1): T[] {
  const target = index + direction;
  if (target < 0 || target >= list.length) return list;

  const next = [...list];
  const [lifted] = next.splice(index, 1);
  if (lifted !== undefined) next.splice(target, 0, lifted);

  return next;
}

export function useRooms() {
  const query = useQuery(roomQueries.list());

  return { rooms: query.data ?? [], isPending: query.isPending };
}

export function useRoomActions({ onDevicesMoved }: { onDevicesMoved?: () => void } = {}) {
  const queryClient = useQueryClient();

  const refreshRooms = () => queryClient.invalidateQueries({ queryKey: roomQueries.all() });

  const create = useMutation({
    mutationFn: (name: string) => createRoom({ name }),
    onSuccess: () => void refreshRooms(),
    onError: (error) => toast.error(reason(error, "Could not create the room")),
  });

  const rename = useMutation({
    mutationFn: ({ id, name }: { id: string; name: string }) => updateRoom(id, { name }),
    onSuccess: (updated) => {
      queryClient.setQueryData<RoomDto[]>(roomQueries.list().queryKey, (current) =>
        current?.map((room) => (room.id === updated.id ? updated : room)),
      );
    },
    onError: (error) => {
      toast.error(reason(error, "Could not rename the room"));
      void refreshRooms();
    },
  });

  const reorder = useMutation({
    mutationFn: async (ordered: RoomDto[]) => {
      await Promise.all(
        ordered.map((room, position) =>
          room.position === position ? undefined : updateRoom(room.id, { position }),
        ),
      );
    },
    onMutate: async (ordered) => {
      await queryClient.cancelQueries({ queryKey: roomQueries.all() });

      const previous = queryClient.getQueryData<RoomDto[]>(roomQueries.list().queryKey);

      queryClient.setQueryData<RoomDto[]>(
        roomQueries.list().queryKey,
        ordered.map((room, position) => ({ ...room, position })),
      );

      return { previous };
    },
    onError: (error, _ordered, context) => {
      if (context?.previous) {
        queryClient.setQueryData(roomQueries.list().queryKey, context.previous);
      }

      toast.error(reason(error, "Could not reorder the rooms"));
    },
    onSettled: () => void refreshRooms(),
  });

  const remove = useMutation({
    mutationFn: (id: string) => deleteRoom(id),
    onSuccess: () => {
      void refreshRooms();
      onDevicesMoved?.();
    },
    onError: (error) => toast.error(reason(error, "Could not delete the room")),
  });

  return { create, rename, reorder, remove };
}
