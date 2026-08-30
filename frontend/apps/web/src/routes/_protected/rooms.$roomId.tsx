import { createFileRoute } from "@tanstack/react-router";

import { RoomPage } from "@/modules/devices/pages/room-page";

export const Route = createFileRoute("/_protected/rooms/$roomId")({
  // Which device's controls are open: a search param, so the room stays mounted underneath.
  validateSearch: (search: Record<string, unknown>): { device?: string } =>
    typeof search.device === "string" && search.device !== "" ? { device: search.device } : {},
  component: RoomPage,
});
