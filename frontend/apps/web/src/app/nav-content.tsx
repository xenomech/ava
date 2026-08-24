import {
  Button,
  Menu,
  MenuContent,
  MenuItem,
  MenuSeparator,
  MenuTrigger,
  cn,
} from "@ava/ui";
import { isOn, type RoomDto } from "@ava/contracts";
import { Link } from "@tanstack/react-router";
import { ChevronDownIcon, ChevronUpIcon, MoreHorizontalIcon, Trash2Icon } from "lucide-react";
import type { CSSProperties } from "react";

import { useSession } from "@/modules/auth";
import { deviceColor } from "@/modules/devices/components/device-stage";
import { useDevices } from "@/modules/devices/use-devices";
import { NewRoom, moved, useRoomActions, useRooms } from "@/modules/rooms";

/**
 * The room rail from the settled desktop design: rooms are the navigation, and
 * each one carries its own name plus how much of it is lit. The foot holds the
 * account rather than the prototype's pair of text links, because settings,
 * hubs and people all live behind it now.
 *
 * Reordering and deleting a room sit on its row here rather than in the room's
 * own heading. Both are about a room's place in this list, and putting them on
 * the page you are looking at meant a permanent row of buttons — one of them
 * destructive — crowding the title on a phone.
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
          {rooms.map((room, at) => (
            <RoomRow
              key={room.id}
              room={room}
              deviceCount={devices.filter((device) => device.room_id === room.id).length}
              lit={devices.filter((device) => device.room_id === room.id && isOn(device))}
              isFirst={at === 0}
              isLast={at === rooms.length - 1}
              onNavigate={onNavigate}
              onMove={(direction) => actions.reorder.mutate(moved(rooms, at, direction))}
              onDelete={() => actions.remove.mutate(room.id)}
            />
          ))}

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

function RoomRow({
  room,
  deviceCount,
  lit,
  isFirst,
  isLast,
  onNavigate,
  onMove,
  onDelete,
}: {
  room: RoomDto;
  deviceCount: number;
  lit: { id: string }[];
  isFirst: boolean;
  isLast: boolean;
  onNavigate?: () => void;
  onMove: (direction: -1 | 1) => void;
  onDelete: () => void;
}) {
  const first = lit[0] as Parameters<typeof deviceColor>[0] | undefined;

  return (
    /* The trigger sits over the link rather than inside it: an anchor may not
       contain a button, and nesting them breaks both. */
    <div className="group/room relative">
      <Link
        to="/rooms/$roomId"
        params={{ roomId: room.id }}
        onClick={onNavigate}
        style={{ "--lit": first ? deviceColor(first) : "var(--color-off)" } as CSSProperties}
        className={cn(
          "grid min-h-[56px] content-center gap-[3px] rounded-[14px] p-3 pr-12",
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
          {lit.length} of {deviceCount} on
        </span>
      </Link>

      <Menu>
        <MenuTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            aria-label={`Actions for ${room.name}`}
            className={cn(
              "absolute right-1 top-1/2 -translate-y-1/2 text-muted",
              "[@media(hover:hover)]:size-9 [@media(hover:hover)]:opacity-0",
              "group-hover/room:opacity-100 focus-visible:opacity-100",
              "data-[state=open]:opacity-100",
            )}
          >
            <MoreHorizontalIcon className="size-4" aria-hidden />
          </Button>
        </MenuTrigger>

        <MenuContent align="end">
          <MenuItem disabled={isFirst} onSelect={() => onMove(-1)}>
            <ChevronUpIcon aria-hidden />
            Move up
          </MenuItem>
          <MenuItem disabled={isLast} onSelect={() => onMove(1)}>
            <ChevronDownIcon aria-hidden />
            Move down
          </MenuItem>

          <MenuSeparator />

          <MenuItem tone="danger" onSelect={onDelete}>
            <Trash2Icon aria-hidden />
            {deviceCount === 0 ? "Delete room" : `Delete, freeing ${deviceCount}`}
          </MenuItem>
        </MenuContent>
      </Menu>
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
