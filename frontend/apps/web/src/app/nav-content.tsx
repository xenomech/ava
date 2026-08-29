import {
  Button,
  Confirm,
  Menu,
  MenuContent,
  MenuItem,
  MenuSeparator,
  MenuTrigger,
  cn,
} from "@ava/ui";
import { isOn, type DeviceDto, type RoomDto } from "@ava/contracts";
import { Link } from "@tanstack/react-router";
import { ChevronDownIcon, ChevronUpIcon, MoreHorizontalIcon, Trash2Icon } from "lucide-react";
import { useMemo, useState, type CSSProperties } from "react";

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

  /* One dialog for the whole list rather than one per row: only ever one room
     is being deleted, and mounting a portal per room to say so is waste. */
  const [doomed, setDoomed] = useState<RoomDto | null>(null);

  /* One pass over the devices for the whole rail, not two filters per room. */
  const byRoom = useMemo(() => {
    const map = new Map<string, { count: number; lit: DeviceDto[] }>();

    for (const device of devices) {
      if (!device.room_id) continue;

      let entry = map.get(device.room_id);
      if (!entry) {
        entry = { count: 0, lit: [] };
        map.set(device.room_id, entry);
      }

      entry.count += 1;
      if (isOn(device)) entry.lit.push(device);
    }

    return map;
  }, [devices]);

  const inDoomed = doomed ? (byRoom.get(doomed.id)?.count ?? 0) : 0;

  return (
    <div className="grid h-full grid-rows-[auto_minmax(0,1fr)_auto]">
      <Confirm
        open={doomed !== null}
        onOpenChange={(open) => !open && setDoomed(null)}
        title={`Delete ${doomed?.name ?? "this room"}?`}
        description={
          inDoomed === 0
            ? "The room goes; nothing else changes."
            : `Its ${inDoomed === 1 ? "device stays" : `${inDoomed} devices stay`} set up, but ${inDoomed === 1 ? "it moves" : "they move"} out of any room, and any scenes saved here are lost.`
        }
        confirmLabel="Delete room"
        tone="danger"
        onConfirm={() => {
          if (doomed) actions.remove.mutate(doomed.id);
          setDoomed(null);
        }}
      />
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
              deviceCount={byRoom.get(room.id)?.count ?? 0}
              lit={byRoom.get(room.id)?.lit ?? NO_DEVICES}
              at={at}
              count={rooms.length}
              onNavigate={onNavigate}
              onMove={(direction) => actions.reorder.mutate(moved(rooms, at, direction))}
              onDelete={() => setDoomed(room)}
            />
          ))}

          {rooms.length === 0 ? (
            <p className="px-3 py-2 text-small text-subtle">No rooms yet.</p>
          ) : null}
        </div>

        <div className="pt-3">
          <NewRoom
            onCreate={(name) => actions.create.mutate(name)}
            busy={actions.create.isPending}
          />
        </div>
      </nav>

      <UserPill onNavigate={onNavigate} />
    </div>
  );
}

const NO_DEVICES: DeviceDto[] = [];

function RoomRow({
  room,
  deviceCount,
  lit,
  at,
  count,
  onNavigate,
  onMove,
  onDelete,
}: {
  room: RoomDto;
  deviceCount: number;
  lit: DeviceDto[];
  /** The row's place in the rail, so the menu knows which moves exist. */
  at: number;
  count: number;
  onNavigate?: () => void;
  onMove: (direction: -1 | 1) => void;
  onDelete: () => void;
}) {
  const first = lit[0];

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
          <MenuItem disabled={at === 0} onSelect={() => onMove(-1)}>
            <ChevronUpIcon aria-hidden />
            Move up
          </MenuItem>
          <MenuItem disabled={at === count - 1} onSelect={() => onMove(1)}>
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
