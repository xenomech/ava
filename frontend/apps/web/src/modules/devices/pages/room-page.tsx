import { cn } from "@ava/ui";
import {
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  emitsLight,
  isOn,
  numberOf,
  supports,
  type DeviceDto,
} from "@ava/contracts";
import { Link, useParams } from "@tanstack/react-router";
import { useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";

import { RoomHeading, rememberRoom, useRoomActions, useRooms } from "@/modules/rooms";
import { Loader } from "@/shared/components/loader";
import { parseColor, warmth } from "@/shared/lib/color";
import { kelvinToCss } from "@/shared/lib/kelvin";
import { Missing } from "../components/empty-state";
import { LightSweep } from "../components/light-sweep";
import { RoomSwitch } from "../components/room-switch";
import { deviceColor, deviceKind, deviceLabel } from "../components/device-stage";
import { deviceQueries } from "../queries";
import { useDevices, useRoomPower } from "../use-devices";

/** The room's own light, averaged, for the switch and the sweep to borrow. */
const DEFAULT_KELVIN = 2700;

/**
 * A room, as one surface: every device at once, a switch in the middle, and the
 * individual devices reduced to a strip you can swipe rather than a grid you
 * have to read.
 */
export function RoomPage() {
  const { roomId } = useParams({ from: "/_protected/rooms/$roomId" });
  const { rooms, isPending: roomsPending } = useRooms();
  const { devices, isPending } = useDevices();
  const queryClient = useQueryClient();
  const setRoomPower = useRoomPower();

  /* The sweep is fire-and-forget: `play` remounts it so a second flick
     restarts the animation instead of being swallowed mid-flight. */
  const [sweep, setSweep] = useState({ play: 0, direction: "on" as "on" | "off" });

  const actions = useRoomActions({
    onDevicesMoved: () => void queryClient.invalidateQueries({ queryKey: deviceQueries.all() }),
  });

  /* So `/` can come back here next time. Writing it on view rather than on
     click means a bookmark or a shared link counts too. */
  useEffect(() => rememberRoom(roomId), [roomId]);

  if (isPending || roomsPending) return <Loader label="Loading room" />;

  const room = rooms.find((entry) => entry.id === roomId);

  if (!room) {
    return (
      <Missing
        title="That room is gone"
        detail="It may have been deleted from another device."
      />
    );
  }

  const inRoom = devices.filter((device) => device.room_id === room.id);
  const on = inRoom.filter(isOn).length;
  const at = rooms.findIndex((entry) => entry.id === room.id);
  const switchable = inRoom.filter(
    (device) => device.status !== "offline" && supports(device, TRAIT_POWER),
  );

  const kelvin = roomKelvin(inRoom);
  const lit = kelvinToCss(kelvin);
  const palette = roomPalette(inRoom, lit);

  const flick = (next: boolean) => {
    setSweep((current) => ({ play: current.play + 1, direction: next ? "on" : "off" }));
    void setRoomPower(inRoom, next);
  };

  return (
    <div className="relative grid h-full grid-rows-[auto_minmax(0,1fr)_auto] overflow-hidden">
      <LightSweep colors={palette} direction={sweep.direction} play={sweep.play} />

      <header className="z-raised flex items-start justify-between gap-4 p-5 pt-16 sm:p-6 md:pt-6">
        <div className="min-w-0">
          <span className="font-mono text-caption uppercase tracking-caps text-subtle">Room</span>
          <div className="mt-1 flex min-w-0 items-center gap-2">
            <RoomHeading
              room={room}
              deviceCount={inRoom.length}
              isFirst={at === 0}
              isLast={at === rooms.length - 1}
              onRename={(name) => actions.rename.mutate({ id: room.id, name })}
              onMove={(direction) => actions.reorder.mutate(moved(rooms, at, direction))}
              onDelete={() => actions.remove.mutate(room.id)}
            />
          </div>
        </div>
      </header>

      {inRoom.length === 0 ? (
        <div className="z-raised grid place-items-center p-6">
          <p className="text-small text-muted">No devices in here yet.</p>
        </div>
      ) : (
        <>
          <div className="z-raised grid min-h-0 place-items-center px-5 sm:px-6">
            <RoomSwitch
              on={on > 0}
              disabled={switchable.length === 0}
              color={lit}
              label={`${room.name} lights`}
              onFlick={flick}
            />
          </div>

          <footer className="z-raised grid gap-4 p-5 pt-0 sm:p-6 sm:pt-0">
            <div className="flex items-end justify-between gap-4">
              <p className="flex items-baseline gap-2">
                {on === 0 ? (
                  <b className="text-hero font-semibold text-subtle">Off</b>
                ) : (
                  <>
                    <b className="text-hero font-semibold tabular">{on}</b>
                    <span className="font-mono text-small text-subtle tabular">
                      of {inRoom.length} on
                    </span>
                  </>
                )}
              </p>

              <p className="text-right font-mono text-caption leading-relaxed text-subtle">
                flick the switch
                <br />
                up on · down off
              </p>
            </div>

            <div
              className="-mx-5 flex snap-x snap-mandatory gap-2.5 overflow-x-auto px-5 pb-1 sm:-mx-6 sm:px-6"
              aria-label={`Devices in ${room.name}`}
            >
              {inRoom.map((device) => (
                <DeviceChip key={device.id} device={device} />
              ))}
            </div>
          </footer>
        </>
      )}
    </div>
  );
}

/** A device reduced to what you need to pick it out: a lit dot and a name. */
function DeviceChip({ device }: { device: DeviceDto }) {
  const live = isOn(device) && device.status !== "offline";

  return (
    <Link
      to="/devices/$deviceId"
      params={{ deviceId: device.id }}
      className={cn(
        "grid w-[168px] shrink-0 snap-start content-start gap-1.5 rounded-lg p-3",
        "border border-border bg-surface transition-colors duration-150 ease-out",
        "hover:border-border-strong focus-visible:outline-none focus-visible:ring-2",
        "focus-visible:ring-fg",
        device.status === "offline" && "opacity-55",
      )}
    >
      <span className="flex items-center gap-2">
        <span
          aria-hidden
          className={cn("size-1.5 shrink-0 rounded-full", !live && "bg-border-strong")}
          style={live ? { background: deviceColor(device) } : undefined}
        />
        <span className="truncate text-small font-semibold">{device.name}</span>
      </span>
      <span className="font-mono text-caption text-subtle tabular">{deviceLabel(device)}</span>
    </Link>
  );
}

/**
 * The colours the room will be holding once the switch settles, warmest first.
 *
 * Every light in the room, not only the ones currently on — flicking up turns
 * all of them on, so the sweep should be showing what is about to be true.
 * Plugs and fans are left out: they have a power state but no colour, and
 * including them would wash the ramp towards whatever `deviceColor` falls back
 * to rather than towards anything the room actually emits.
 */
function roomPalette(devices: DeviceDto[], fallback: string): string[] {
  const lights = devices
    .filter((device) => emitsLight(deviceKind(device)))
    .map((device) => deviceColor(device));

  const seen = new Set<string>();
  const unique = lights.filter((color) => !seen.has(color) && seen.add(color));

  if (unique.length === 0) return [fallback];

  return unique
    .map((color) => ({ color, rgb: parseColor(color) }))
    .filter((entry) => entry.rgb !== null)
    .sort((a, b) => warmth(b.rgb ?? [0, 0, 0]) - warmth(a.rgb ?? [0, 0, 0]))
    .map((entry) => entry.color);
}

function roomKelvin(devices: DeviceDto[]): number {
  const temps = devices
    .map((device) => numberOf(device, TRAIT_COLOR_TEMP))
    .filter((value): value is number => typeof value === "number");

  if (temps.length === 0) return DEFAULT_KELVIN;

  return Math.round(temps.reduce((sum, value) => sum + value, 0) / temps.length);
}

function moved<T>(list: T[], index: number, direction: -1 | 1): T[] {
  const target = index + direction;
  if (target < 0 || target >= list.length) return list;

  const next = [...list];
  const [lifted] = next.splice(index, 1);
  if (lifted !== undefined) next.splice(target, 0, lifted);

  return next;
}
