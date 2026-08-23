import { isOn, type DeviceDto, type RoomDto } from "@ava/contracts";
import { useMemo } from "react";

import { useRooms } from "@/modules/rooms";

export const NO_ROOM = "no-room";

export type RoomGroup = {
  /** Room id, or NO_ROOM for devices nobody has filed yet. */
  key: string;
  name: string;
  room: RoomDto | undefined;
  devices: DeviceDto[];
  on: number;
};

// The console is room-first: a fixture is always shown inside the room you are
// standing in, and flicking never leaves it. That makes the grouping — not the
// flat device list — the structure every part of the page reads from.
export function useRoomGroups(devices: DeviceDto[]): RoomGroup[] {
  const { rooms } = useRooms();

  return useMemo(() => {
    const known = new Set(rooms.map((room) => room.id));

    // Bucket once rather than filtering the whole device list per room, which
    // would be rooms x devices.
    const byRoom = new Map<string, DeviceDto[]>();

    for (const device of devices) {
      const key = device.room_id && known.has(device.room_id) ? device.room_id : NO_ROOM;
      const bucket = byRoom.get(key);

      if (bucket) bucket.push(device);
      else byRoom.set(key, [device]);
    }

    const grouped = rooms.map((room) => group(room.id, room.name, room, byRoom));

    if (byRoom.has(NO_ROOM)) {
      grouped.push(group(NO_ROOM, "Unassigned", undefined, byRoom));
    }

    // A room with nothing in it has nothing to control, so it never becomes the
    // console's subject. It stays visible on the Rooms page.
    return grouped.filter((entry) => entry.devices.length > 0);
  }, [rooms, devices]);
}

function group(
  key: string,
  name: string,
  room: RoomDto | undefined,
  byRoom: Map<string, DeviceDto[]>,
): RoomGroup {
  const inRoom = byRoom.get(key) ?? [];
  let on = 0;

  for (const device of inRoom) {
    if (isOn(device)) on += 1;
  }

  return { key, name, room, devices: inRoom, on };
}
