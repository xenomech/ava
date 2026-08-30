import { Button, Chip, Device, DeviceHalo, Switch, cn } from "@ava/ui";
import { TRAIT_POWER, isOn, supports, type DeviceDto } from "@ava/contracts";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";

import { hubQueries } from "@/modules/hub";
import { NewRoom, RoomHeading, useRoomActions, useRooms } from "@/modules/rooms";
import { Loader } from "@/shared/components/loader";
import { deviceColor, deviceKind, deviceLevel } from "../components/device-stage";
import { NoDevices } from "../components/empty-state";
import { deviceQueries } from "../queries";
import { useDeviceControl, useDevices, useRoomPower } from "../use-devices";

const UNASSIGNED = "No room";

export function RoomsPage() {
  const { devices, isPending } = useDevices();
  const hubs = useQuery(hubQueries.list());
  const { rooms: defined } = useRooms();
  const queryClient = useQueryClient();
  const actions = useRoomActions({
    onDevicesMoved: () => void queryClient.invalidateQueries({ queryKey: deviceQueries.all() }),
  });
  const control = useDeviceControl();
  const setRoom = useRoomPower();

  if (isPending) return <Loader label="Loading devices" />;

  if (devices.length === 0) return <NoDevices hasHub={(hubs.data ?? []).length > 0} />;

  const known = new Set(defined.map((room) => room.id));
  const unassigned = devices.filter((device) => !device.room_id || !known.has(device.room_id));
  const groups = [
    ...defined.map((room) => ({
      key: room.id,
      name: room.name,
      room,
      devices: devices.filter((device) => device.room_id === room.id),
    })),
    ...(unassigned.length > 0
      ? [{ key: "unassigned", name: UNASSIGNED, room: undefined, devices: unassigned }]
      : []),
  ];

  return (
    <div className="grid gap-8 p-5 sm:p-6">
      <header className="flex items-start justify-between gap-4">
        <div>
          <h1 className="text-display font-semibold">Rooms</h1>
          <p className="mt-1 text-small text-muted">
            {groups.length} {groups.length === 1 ? "room" : "rooms"} · {devices.length}{" "}
            {devices.length === 1 ? "device" : "devices"}
          </p>
        </div>

        <NewRoom onCreate={(name) => actions.create.mutate(name)} busy={actions.create.isPending} />
      </header>

      {groups.map((group, index) => {
        const inRoom = group.devices;
        const on = inRoom.filter((device) => isOn(device)).length;

        return (
          <section key={group.key} className="grid gap-3">
            <div className="flex items-center justify-between gap-4">
              <div className="flex min-w-0 items-center gap-3">
                {group.room ? (
                  <RoomHeading
                    room={group.room}
                    deviceCount={inRoom.length}
                    isFirst={index === 0}
                    isLast={index === defined.length - 1}
                    onRename={(name) => actions.rename.mutate({ id: group.room.id, name })}
                    onMove={(direction) => actions.reorder.mutate(moved(defined, index, direction))}
                    onDelete={() => actions.remove.mutate(group.room.id)}
                  />
                ) : (
                  <h2 className="text-title font-semibold text-muted">{group.name}</h2>
                )}
                <span className="font-mono text-caption text-subtle tabular">
                  {inRoom.length === 0 ? "empty" : `${on} of ${inRoom.length} on`}
                </span>
              </div>

              <RoomPower
                room={group.name}
                devices={inRoom}
                anyOn={on > 0}
                onSet={(next) => void setRoom(inRoom, next)}
              />
            </div>

            <div className="grid grid-cols-[repeat(auto-fill,minmax(210px,1fr))] gap-3">
              {inRoom.length === 0 ? (
                <p className="text-small text-muted">No devices in here yet.</p>
              ) : null}
              {inRoom.map((device) => (
                <DeviceCard
                  key={device.id}
                  device={device}
                  disabled={device.status === "offline"}
                  onToggle={() => control(device, TRAIT_POWER, !isOn(device))}
                />
              ))}
            </div>
          </section>
        );
      })}
    </div>
  );
}

function DeviceCard({
  device,
  disabled,
  onToggle,
}: {
  device: DeviceDto;
  disabled: boolean;
  onToggle: () => void;
}) {
  const level = deviceLevel(device);
  const color = deviceColor(device);

  return (
    <article
      className={cn(
        "grid gap-3 rounded-lg border border-border bg-surface p-4",
        "transition-colors duration-150 ease-out",
        isOn(device) && "border-border-strong",
        device.status === "offline" && "opacity-60",
      )}
      style={{ "--level": level, "--lit": color } as React.CSSProperties}
    >
      <Link
        to="/"
        search={{ device: device.id }}
        className="relative grid h-32 place-items-center rounded-sm"
        aria-label={`Open ${device.name}`}
      >
        <DeviceHalo className="w-3/4" />
        <Device kind={deviceKind(device)} level={level} color={color} className="h-[86%]" />
      </Link>

      <div className="flex items-end justify-between gap-3">
        <span className="min-w-0">
          <span className="block truncate text-body font-semibold">{device.name}</span>
          <span className="block font-mono text-caption text-subtle tabular">
            {device.status === "offline" ? "Offline" : level > 0 ? `${level}%` : "Off"}
          </span>
        </span>

        {device.status === "offline" ? (
          <Chip tone="warning">Offline</Chip>
        ) : (
          <Switch
            checked={isOn(device)}
            disabled={disabled}
            onCheckedChange={onToggle}
            aria-label={device.name}
          />
        )}
      </div>
    </article>
  );
}

function RoomPower({
  room,
  devices,
  anyOn,
  onSet,
}: {
  room: string;
  devices: DeviceDto[];
  anyOn: boolean;
  onSet: (on: boolean) => void;
}) {
  const switchable = devices.filter(
    (device) => device.status !== "offline" && supports(device, TRAIT_POWER),
  );

  if (switchable.length < 2) return null;

  return (
    <div className="flex shrink-0 gap-1.5">
      <Button
        variant="ghost"
        size="sm"
        disabled={!anyOn}
        onClick={() => onSet(false)}
        aria-label={`Turn everything off in ${room}`}
      >
        All off
      </Button>
      <Button
        variant="secondary"
        size="sm"
        onClick={() => onSet(true)}
        aria-label={`Turn everything on in ${room}`}
      >
        All on
      </Button>
    </div>
  );
}

function moved<T>(list: T[], index: number, direction: -1 | 1): T[] {
  const target = index + direction;
  if (target < 0 || target >= list.length) return list;

  const next = [...list];
  const [lifted] = next.splice(index, 1);
  if (lifted !== undefined) next.splice(target, 0, lifted);

  return next;
}
