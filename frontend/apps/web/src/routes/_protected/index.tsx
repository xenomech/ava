import { createFileRoute, redirect } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";

import { NoDevices, NoRooms, deviceQueries } from "@/modules/devices";
import { hubQueries } from "@/modules/hub";
import { recallRoom, roomQueries } from "@/modules/rooms";
import { Loader } from "@/shared/components/loader";

/** `/` is not a page: it redirects to the last room, and only renders when there is none. */
export const Route = createFileRoute("/_protected/")({
  beforeLoad: async ({ context }) => {
    const rooms = await context.queryClient.ensureQueryData(roomQueries.list());

    if (rooms.length === 0) return;

    const last = recallRoom();
    // A remembered id from another account is not in this list, so the check keeps it private.
    const target = rooms.find((room) => room.id === last) ?? rooms[0];

    if (target) {
      throw redirect({ to: "/rooms/$roomId", params: { roomId: target.id }, replace: true });
    }
  },
  component: Landing,
});

function Landing() {
  const devices = useQuery(deviceQueries.list());
  const hubs = useQuery(hubQueries.list());

  if (devices.isPending) return <Loader label="Loading home" />;

  const loose = devices.data ?? [];

  if (loose.length === 0) return <NoDevices hasHub={(hubs.data ?? []).length > 0} />;

  return <NoRooms devices={loose} />;
}
