import { createFileRoute } from "@tanstack/react-router";

import { RoomPage } from "@/modules/devices/pages/room-page";

export const Route = createFileRoute("/_protected/rooms/$roomId")({
  component: RoomPage,
});
