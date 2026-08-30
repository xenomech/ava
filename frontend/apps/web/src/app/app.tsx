import { Link } from "@tanstack/react-router";

import { UserPill, useSession } from "@/modules/auth";
import { CommandPalette, useDeviceEvents } from "@/modules/devices";
import { useHubEvents } from "@/modules/hub";
import { RoomRail } from "@/modules/rooms";
import { AvaSocketProvider } from "@/shared/realtime";
import { Shell } from "./shell";

/** The composition root: the one place that knows which modules fill the shell's slots. */
export function AppShell() {
  return (
    <AvaSocketProvider>
      <RealtimeSync />
      <Shell nav={<AppNav />} palette={<CommandPalette />} />
    </AvaSocketProvider>
  );
}

// Each module keeps its own cache in step with the socket; mounting the hooks is all the wiring.
function RealtimeSync() {
  useDeviceEvents();
  useHubEvents();

  return null;
}

function AppNav() {
  const { tenant } = useSession();

  return (
    <div className="grid h-full grid-rows-[auto_minmax(0,1fr)_auto]">
      <Link to="/" className="flex items-baseline gap-2.5 px-3 pb-8">
        <b className="text-[1.0625rem] font-semibold tracking-tight">Ava</b>
        <span className="truncate font-mono text-micro text-subtle">{tenant?.name ?? ""}</span>
      </Link>

      <RoomRail />

      <UserPill />
    </div>
  );
}
