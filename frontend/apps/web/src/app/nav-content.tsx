import { cn } from "@ava/ui";
import { isOn } from "@ava/contracts";
import { Link } from "@tanstack/react-router";
import type { CSSProperties } from "react";

import { useSession } from "@/modules/auth";
import { deviceColor } from "@/modules/devices/components/device-stage";
import { useDevices } from "@/modules/devices/use-devices";
import { NewRoom, useRoomActions, useRooms } from "@/modules/rooms";

/**
 * The room rail from the settled desktop design: rooms are the navigation, and
 * each one carries its own name plus how much of it is lit. The foot holds the
 * account rather than the prototype's pair of text links, because settings,
 * hubs and people all live behind it now.
 *
 * Shared by the desktop rail and the phone drawer so the two cannot drift.
 */
export function NavContent({ onNavigate }: { onNavigate?: () => void }) {
  const { rooms } = useRooms();
  const { devices } = useDevices();
  const { tenant } = useSession();
  const actions = useRoomActions();

  return (
    <div className="grid h-full grid-rows-[auto_minmax(0,1fr)_auto]">
      <Link to="/" onClick={onNavigate} className="flex items-baseline gap-2.5 px-3 pb-8">
        <b className="text-[1.0625rem] font-semibold tracking-tight">Ava</b>
        <span className="truncate font-mono text-micro text-subtle">{tenant?.name ?? ""}</span>
      </Link>

      <nav aria-label="Rooms" className="min-h-0 overflow-y-auto">
        <div className="grid gap-2">
          {rooms.map((room) => {
            const inRoom = devices.filter((device) => device.room_id === room.id);
            const lit = inRoom.filter(isOn);

            return (
              <Link
                key={room.id}
                to="/rooms/$roomId"
                params={{ roomId: room.id }}
                onClick={onNavigate}
                style={
                  { "--lit": lit[0] ? deviceColor(lit[0]) : "var(--color-off)" } as CSSProperties
                }
                className={cn(
                  "grid min-h-[56px] content-center gap-[3px] rounded-[14px] p-3",
                  "text-muted transition-colors duration-[180ms] ease-out-soft",
                  "hover:bg-raised",
                  "aria-[current=page]:bg-raised aria-[current=page]:text-fg",
                )}
              >
                <b className="truncate text-[0.9375rem] font-medium">{room.name}</b>
                <span className="flex items-center gap-1.5 font-mono text-micro text-subtle">
                  <span
                    aria-hidden
                    className="size-[5px] shrink-0 rounded-full"
                    style={{
                      background: lit.length > 0 ? "var(--lit)" : "var(--color-off)",
                      boxShadow: lit.length > 0 ? "0 0 6px var(--lit)" : "none",
                    }}
                  />
                  {lit.length} of {inRoom.length} on
                </span>
              </Link>
            );
          })}

          {rooms.length === 0 ? (
            <p className="px-3 py-2 text-small text-subtle">No rooms yet.</p>
          ) : null}
        </div>

        <div className="pt-3">
          <NewRoom onCreate={(name) => actions.create.mutate(name)} busy={actions.create.isPending} />
        </div>
      </nav>

      <UserPill onNavigate={onNavigate} />
    </div>
  );
}

// Where the prototype had "Add a hub" and "Settings", one target now stands for
// the account and everything configurable behind it.
function UserPill({ onNavigate }: { onNavigate?: () => void }) {
  const { user, tenant } = useSession();
  const initial = (user?.name ?? user?.email ?? "?").trim().slice(0, 1).toUpperCase();

  return (
    <Link
      to="/settings"
      onClick={onNavigate}
      className={cn(
        "mt-4 flex min-h-[56px] items-center gap-3 rounded-[14px] p-2.5",
        "text-muted transition-colors duration-[180ms] ease-out-soft",
        "hover:bg-raised hover:text-fg",
        "aria-[current=page]:bg-raised aria-[current=page]:text-fg",
      )}
    >
      <span
        aria-hidden
        className="grid size-9 shrink-0 place-items-center rounded-full bg-raised text-small font-semibold text-fg"
      >
        {initial}
      </span>

      <span className="min-w-0 flex-1">
        <span className="block truncate text-[0.8125rem] font-medium text-fg">
          {user?.name ?? "Account"}
        </span>
        <span className="block truncate font-mono text-micro text-subtle">
          {tenant?.name ?? "ava"}
        </span>
      </span>
    </Link>
  );
}
