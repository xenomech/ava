import {
  roomDto,
  type CreateRoomRequest,
  type RoomDto,
  type UpdateRoomRequest,
} from "@ava/contracts";
import { z } from "zod";

import { request } from "@/config/http/request";

export function listRooms(signal?: AbortSignal): Promise<RoomDto[]> {
  return request({ url: "/rooms", schema: z.array(roomDto), signal });
}

export function createRoom(body: CreateRoomRequest): Promise<RoomDto> {
  return request({ url: "/rooms", method: "post", body, schema: roomDto });
}

export function updateRoom(roomID: string, body: UpdateRoomRequest): Promise<RoomDto> {
  return request({ url: `/rooms/${roomID}`, method: "patch", body, schema: roomDto });
}

export function deleteRoom(roomID: string): Promise<void> {
  return request({ url: `/rooms/${roomID}`, method: "delete" });
}
