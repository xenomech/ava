import { Device, DeviceHalo, cn } from "@ava/ui";
import {
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  emitsLight,
  isOn,
  numberOf,
  supports,
  type DeviceDto,
  type SceneDto,
} from "@ava/contracts";
import { useNavigate, useParams, useSearch } from "@tanstack/react-router";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback, useEffect, useMemo, useState } from "react";

import { HubPill, hubQueries } from "@/modules/hub";
import { RoomHeading, rememberRoom, useRoomActions, useRooms } from "@/modules/rooms";
import {
  SceneRow,
  sceneColor,
  scenePreview,
  useApplyScene,
  useArmedScene,
  useScenes,
} from "@/modules/scenes";
import { Loader } from "@/shared/components/loader";
import { useMediaQuery } from "@/shared/hooks/use-media-query";
import { parseColor, warmth } from "@/shared/lib/color";
import { kelvinToCss } from "@/shared/lib/kelvin";
import { deviceColor, deviceKind, deviceLevel } from "../lib/device-view";
import { DeviceDrawer, DevicePanel } from "../components/device-sheet";
import { BESIDE, ROOM_HEIGHT } from "../constants";
import { DeviceStrip } from "../components/device-strip";
import { Missing } from "../components/empty-state";
import { LightSweep } from "../components/light-sweep";
import { RoomSwitch } from "../components/room-switch";
import { deviceQueries } from "../queries";
import { useDevices, useRoomPower } from "../hooks/use-devices";

/** The room's own light, averaged, for the switch and the sweep to borrow. */
const DEFAULT_KELVIN = 2700;

/** A room as one surface: nothing navigates away, and the selection lives in the URL. */
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

  const { scenes, isPending: scenesPending } = useScenes(roomId);
  const { armed, arm } = useArmedScene(roomId, scenes);
  const applyScene = useApplyScene();

  // `play` remounts the sweep, so a second flick restarts it instead of being swallowed.
  const [sweep, setSweep] = useState({ play: 0, direction: "on" as "on" | "off" });
  // The dragged device's in-flight brightness, so the stage lights before the hub answers.
  const [dragging, setDragging] = useState<number | null>(null);

  const actions = useRoomActions({
    onDevicesMoved: () => void queryClient.invalidateQueries({ queryKey: deviceQueries.all() }),
  });

  // So `/` comes back here; written on view so a bookmark or shared link counts too.
  useEffect(() => rememberRoom(roomId), [roomId]);

  // Memoised because the page re-renders at pointer rate and none of this changes then.
  const inRoom = useMemo(
    () => devices.filter((device) => device.room_id === roomId),
    [devices, roomId],
  );
  const elsewhere = useMemo(
    () => devices.filter((device) => device.room_id !== roomId),
    [devices, roomId],
  );

  const { on, switchable } = useMemo(
    () => ({
      on: inRoom.filter(isOn).length,
      switchable: inRoom.filter(
        (device) => device.status !== "offline" && supports(device, TRAIT_POWER),
      ),
    }),
    [inRoom],
  );

  // Only the hubs answering for this room, since no other hub bears on whether it works.
  const serving = useMemo(() => {
    const hubIds = new Set(inRoom.map((device) => device.hub_id));

    return (hubs.data ?? []).filter((entry) => hubIds.has(entry.id));
  }, [hubs.data, inRoom]);

  // The armed scene's colour, so arming by scroll is visible before anything is sent.
  const ambient = kelvinToCss(roomKelvin(inRoom));
  const lit = useMemo(
    () => (armed ? sceneColor(scenePreview(armed, inRoom), ambient) : ambient),
    [armed, inRoom, ambient],
  );
  const palette = useMemo(() => roomPalette(inRoom, lit), [inRoom, lit]);

  // Stable, so SceneRow's scroll listener is not re-attached per render.
  const onArm = useCallback((scene: SceneDto | null) => arm(scene?.id ?? null), [arm]);

  if (isPending || roomsPending) return <Loader label="Loading room" />;

  const room = rooms.find((entry) => entry.id === roomId);

  if (!room) {
    return (
      <Missing title="That room is gone" detail="It may have been deleted from another device." />
    );
  }

  const selected = inRoom.find((device) => device.id === selectedId);
  const hub = (hubs.data ?? []).find((entry) => entry.id === selected?.hub_id);
  const hubOffline = hub !== undefined && !hub.online;
  const connectivity = hubOffline
    ? ("hub-offline" as const)
    : selected?.status === "offline"
      ? ("device-offline" as const)
      : ("online" as const);

  // Up applies the armed scene, or everything on; down is always just off.
  const flick = (next: boolean) => {
    setSweep((current) => ({ play: current.play + 1, direction: next ? "on" : "off" }));

    if (next && armed) {
      void applyScene(armed);

      return;
    }

    void setRoomPower(inRoom, next);
  };

  const play = (scene: SceneDto | null) => {
    arm(scene?.id ?? null);
    setSweep((current) => ({ play: current.play + 1, direction: "on" }));

    if (scene) void applyScene(scene);
    else void setRoomPower(inRoom, true);
  };

  const close = () => void navigate({ to: "/rooms/$roomId", params: { roomId }, replace: true });

  const strip = (
    <DeviceStrip
      devices={inRoom}
      elsewhere={elsewhere}
      roomId={room.id}
      roomName={room.name}
      selectedId={selected?.id}
      label={`Devices in ${room.name}`}
    />
  );

  return (
    <div className="relative flex h-full overflow-hidden">
      <div
        className={cn(
          "relative grid min-h-0 min-w-0 flex-1 grid-rows-[auto_minmax(0,1fr)_auto] overflow-hidden",
          // Laid across rather than down when the screen is short and wide.
          "landscape-room:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]",
          "landscape-room:grid-rows-[auto_minmax(0,1fr)]",
          // The phone sheet takes the lower half, so the room gives that half up.
          selected && !beside && ROOM_HEIGHT,
        )}
      >
        <LightSweep colors={palette} direction={sweep.direction} play={sweep.play} />

        <header className="z-raised p-5 pt-11 sm:p-6 md:pt-6 landscape-room:col-span-2 landscape-room:pt-4">
          <span className="font-mono text-caption uppercase tracking-caps text-subtle">Room</span>
          <div className="mt-1 flex min-w-0 items-center gap-2">
            <RoomHeading
              room={room}
              onRename={(name) => actions.rename.mutate({ id: room.id, name })}
            />

            {serving.map((entry) => (
              <HubPill key={entry.id} hub={entry} />
            ))}
          </div>
        </header>

        {inRoom.length === 0 ? (
          // An empty room still gets the strip, the only way to put something in it.
          <>
            <div className="z-raised grid place-items-center px-6">
              <p className="max-w-[280px] text-center text-small text-muted">
                Nothing in here yet. Add a device to give this room a switch.
              </p>
            </div>

            <footer className="z-raised p-5 sm:p-6">{strip}</footer>
          </>
        ) : (
          <>
            <div className="z-raised grid min-h-0 place-items-center px-5 sm:px-6">
              {selected ? (
                <DeviceOnStage key={selected.id} device={selected} level={dragging} />
              ) : (
                <div className="grid min-h-0 w-full grid-cols-[minmax(0,1fr)] justify-items-center gap-4">
                  <RoomSwitch
                    on={on > 0}
                    disabled={switchable.length === 0}
                    color={lit}
                    label={`${room.name} lights`}
                    onFlick={flick}
                  />

                  <SceneRow
                    roomId={room.id}
                    roomName={room.name}
                    devices={inRoom}
                    scenes={scenes}
                    scenesReady={!scenesPending}
                    armedId={armed?.id ?? null}
                    onArm={onArm}
                    onApply={play}
                  />
                </div>
              )}
            </div>

            {/* Room mode's furniture, which travels into the sheet on a phone. */}
            {selected && !beside ? null : (
              <footer
                className={cn(
                  "z-raised grid gap-4 p-5 pt-0 sm:p-6 sm:pt-0",
                  "landscape-room:content-center landscape-room:pt-5",
                )}
              >
                {/* First to go on a short screen: the strip already names every device. */}
                <p className="hidden items-baseline gap-2 [@media(min-height:700px)]:flex">
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

                {strip}
              </footer>
            )}
          </>
        )}
      </div>

      {selected ? (
        beside ? (
          <DevicePanel
            device={selected}
            connectivity={connectivity}
            onClose={close}
            onLevelChange={setDragging}
          />
        ) : (
          <DeviceDrawer
            device={selected}
            connectivity={connectivity}
            onClose={close}
            onLevelChange={setDragging}
          >
            {strip}
          </DeviceDrawer>
        )
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

// Every light the room will hold once the switch settles, warmest first, plugs left out.
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
