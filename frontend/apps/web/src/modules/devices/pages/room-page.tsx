import { Device, DeviceHalo, cn } from "@ava/ui";
import {
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  emitsLight,
  isOn,
  numberOf,
  supports,
  type DeviceDto,
} from "@ava/contracts";
import { useNavigate, useParams, useSearch } from "@tanstack/react-router";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";

import { hubQueries } from "@/modules/hub";
import { RoomHeading, rememberRoom, useRoomActions, useRooms } from "@/modules/rooms";
import { Loader } from "@/shared/components/loader";
import { useMediaQuery } from "@/shared/hooks/use-media-query";
import { parseColor, warmth } from "@/shared/lib/color";
import { kelvinToCss } from "@/shared/lib/kelvin";
import { deviceColor, deviceKind, deviceLevel } from "../components/device-stage";
import { DeviceSheet, ROOM_HEIGHT } from "../components/device-sheet";
import { DeviceStrip } from "../components/device-strip";
import { Missing } from "../components/empty-state";
import { LightSweep } from "../components/light-sweep";
import { RoomSwitch } from "../components/room-switch";
import { deviceQueries } from "../queries";
import { useDevices, useRoomPower } from "../use-devices";

/** The room's own light, averaged, for the switch and the sweep to borrow. */
const DEFAULT_KELVIN = 2700;

/** Below this the controls arrive as a sheet; above it, as a column. */
const BESIDE = "(min-width: 1024px)";

/**
 * A room, as one surface.
 *
 * Nothing here navigates away. Picking a device out of the strip swaps the
 * middle of the room from its switch to the device itself and opens that
 * device's controls, but the room, its colours and its strip all stay put. The
 * selection lives in the URL, so it stays addressable and the back button
 * closes it.
 */
export function RoomPage() {
  const { roomId } = useParams({ from: "/_protected/rooms/$roomId" });
  const { device: selectedId } = useSearch({ from: "/_protected/rooms/$roomId" });
  const navigate = useNavigate();

  const { rooms, isPending: roomsPending } = useRooms();
  const { devices, isPending } = useDevices();
  const hubs = useQuery(hubQueries.list());
  const queryClient = useQueryClient();
  const setRoomPower = useRoomPower();
  const beside = useMediaQuery(BESIDE);

  /* The sweep is fire-and-forget: `play` remounts it so a second flick
     restarts the animation instead of being swallowed mid-flight. */
  const [sweep, setSweep] = useState({ play: 0, direction: "on" as "on" | "off" });
  /* The in-flight brightness of whichever device is being dragged, so the
     stage lights up before the hub has answered. */
  const [dragging, setDragging] = useState<number | null>(null);

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
      <Missing title="That room is gone" detail="It may have been deleted from another device." />
    );
  }

  const inRoom = devices.filter((device) => device.room_id === room.id);
  const on = inRoom.filter(isOn).length;
  const at = rooms.findIndex((entry) => entry.id === room.id);
  const switchable = inRoom.filter(
    (device) => device.status !== "offline" && supports(device, TRAIT_POWER),
  );

  const selected = inRoom.find((device) => device.id === selectedId);
  const hub = (hubs.data ?? []).find((entry) => entry.id === selected?.hub_id);
  const hubOffline = hub !== undefined && !hub.online;

  const kelvin = roomKelvin(inRoom);
  const lit = kelvinToCss(kelvin);
  const palette = roomPalette(inRoom, lit);

  const flick = (next: boolean) => {
    setSweep((current) => ({ play: current.play + 1, direction: next ? "on" : "off" }));
    void setRoomPower(inRoom, next);
  };

  const close = () => void navigate({ to: "/rooms/$roomId", params: { roomId }, replace: true });

  const strip = (
    <DeviceStrip
      devices={inRoom}
      roomId={room.id}
      selectedId={selected?.id}
      label={`Devices in ${room.name}`}
    />
  );

  return (
    <div className="relative flex h-full overflow-hidden">
      <div
        className={cn(
          "relative grid min-h-0 min-w-0 flex-1 grid-rows-[auto_minmax(0,1fr)_auto] overflow-hidden",
          /* The phone sheet rests over the lower half, so the room gives up
             that half rather than centring the device behind it. */
          selected && !beside && ROOM_HEIGHT,
        )}
      >
        <LightSweep colors={palette} direction={sweep.direction} play={sweep.play} />

        <header className="z-raised p-5 pt-16 sm:p-6 md:pt-6">
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
        </header>

        {inRoom.length === 0 ? (
          <div className="z-raised grid place-items-center p-6">
            <p className="text-small text-muted">No devices in here yet.</p>
          </div>
        ) : (
          <>
            <div className="z-raised grid min-h-0 place-items-center px-5 sm:px-6">
              {selected ? (
                <DeviceOnStage key={selected.id} device={selected} level={dragging} />
              ) : (
                <RoomSwitch
                  on={on > 0}
                  disabled={switchable.length === 0}
                  color={lit}
                  label={`${room.name} lights`}
                  onFlick={flick}
                />
              )}
            </div>

            {/* Room mode's own furniture. On a phone the sheet takes the lower
                half of the screen, so the count, the hint and the strip all
                travel into it rather than sitting on top of the device. */}
            {selected && !beside ? null : (
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
                    {selected ? (
                      <>
                        tap the card again
                        <br />
                        to leave the device
                      </>
                    ) : (
                      <>
                        flick the switch
                        <br />
                        up on · down off
                      </>
                    )}
                  </p>
                </div>

                {strip}
              </footer>
            )}
          </>
        )}
      </div>

      {selected ? (
        <DeviceSheet
          device={selected}
          offline={selected.status === "offline" || hubOffline}
          hubOffline={hubOffline}
          strip={strip}
          beside={beside}
          onClose={close}
          onLevelChange={setDragging}
        />
      ) : null}
    </div>
  );
}

/** The device itself, standing where the room's switch was. */
function DeviceOnStage({ device, level }: { device: DeviceDto; level: number | null }) {
  const shown = level ?? deviceLevel(device);
  const color = deviceColor(device);

  return (
    <div
      className="relative grid size-full animate-fade-in place-items-center"
      style={{ "--level": shown, "--lit": color } as React.CSSProperties}
    >
      {emitsLight(deviceKind(device)) ? <DeviceHalo className="w-[46%]" /> : null}
      <Device
        kind={deviceKind(device)}
        level={shown}
        color={color}
        className={cn("h-[58%] max-h-[380px]", device.status === "offline" && "opacity-50")}
      />
    </div>
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
