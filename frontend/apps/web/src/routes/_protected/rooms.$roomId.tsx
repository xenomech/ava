import { createFileRoute } from "@tanstack/react-router";

import { RoomPage } from "@/modules/devices/pages/room-page";

export const Route = createFileRoute("/_protected/rooms/$roomId")({
  /* Which device's controls are open. A search param rather than a child route
     because the room stays mounted underneath — this opens a panel, it does not
     navigate anywhere — and it gets a working back button for free. */
  validateSearch: (search: Record<string, unknown>): { device?: string } =>
    typeof search.device === "string" && search.device !== "" ? { device: search.device } : {},
  component: RoomPage,
});
