import { createFileRoute } from "@tanstack/react-router";

import { RoomsPage } from "@/modules/devices/pages/rooms-page";

export const Route = createFileRoute("/_protected/rooms")({
  component: RoomsPage,
});
