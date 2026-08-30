import { createFileRoute, redirect } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";

import { NoDevices, NoRooms } from "@/modules/devices/components/empty-state";
import { deviceQueries } from "@/modules/devices/queries";
import { hubQueries } from "@/modules/hub";
import { recallRoom, roomQueries } from "@/modules/rooms";
import { Loader } from "@/shared/components/loader";

/**
 * `/` is not a page. Rooms are the surface the app is built around, so the
 * root sends you to the one you were last in — or the first one, on a new
 * browser. It only renders anything at all when there is no room to send you
 * to, which is the state a brand new account starts in.
 */
export const Route = createFileRoute("/_protected/")({
  beforeLoad: async ({ context }) => {
    const rooms = await context.queryClient.ensureQueryData(roomQueries.list());

    if (rooms.length === 0) return;

    const last = recallRoom();
    /* A remembered id from another account will not be in this list, so the
       check that keeps the redirect valid is also what keeps it private. */
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
