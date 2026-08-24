import { Button, Chip, Device, DeviceHalo, Switch, cn } from "@ava/ui";
import { TRAIT_POWER, isOn, supports, type DeviceDto } from "@ava/contracts";
import { Link, useParams } from "@tanstack/react-router";
import { useQueryClient } from "@tanstack/react-query";

import { RoomHeading, useRoomActions, useRooms } from "@/modules/rooms";
import { Loader } from "@/shared/components/loader";
import { deviceColor, deviceKind, deviceLevel } from "../components/device-stage";
import { deviceQueries } from "../queries";
import { useDeviceControl, useDevices, useRoomPower } from "../use-devices";

/** One room at a time, reached from the sidebar. */
export function RoomPage() {
  const { roomId } = useParams({ from: "/_protected/rooms/$roomId" });
  const { rooms, isPending: roomsPending } = useRooms();
  const { devices, isPending } = useDevices();
  const queryClient = useQueryClient();
  const control = useDeviceControl();
  const setRoomPower = useRoomPower();

  const actions = useRoomActions({
    onDevicesMoved: () => void queryClient.invalidateQueries({ queryKey: deviceQueries.all() }),
  });

  if (isPending || roomsPending) return <Loader label="Loading room" />;

  const room = rooms.find((entry) => entry.id === roomId);

  if (!room) {
    return (
      <div className="grid min-h-full place-items-center p-6">
        <div className="grid max-w-[360px] justify-items-center gap-3 text-center">
          <h1 className="text-title font-semibold">That room is gone</h1>
          <p className="text-small text-muted">It may have been deleted from another device.</p>
          <Link to="/" className="mt-2">
            <Button>Back home</Button>
          </Link>
        </div>
      </div>
    );
  }

  const inRoom = devices.filter((device) => device.room_id === room.id);
  const on = inRoom.filter(isOn).length;
  const at = rooms.findIndex((entry) => entry.id === room.id);
  const switchable = inRoom.filter(
    (device) => device.status !== "offline" && supports(device, TRAIT_POWER),
  );

  return (
    <div className="mx-auto grid w-full max-w-[1180px] gap-6 p-5 pt-16 sm:p-6 md:pt-6">
      <header className="flex items-center justify-between gap-4">
        <div className="flex min-w-0 items-center gap-3">
          <RoomHeading
            room={room}
            deviceCount={inRoom.length}
            isFirst={at === 0}
            isLast={at === rooms.length - 1}
            onRename={(name) => actions.rename.mutate({ id: room.id, name })}
            onMove={(direction) => actions.reorder.mutate(moved(rooms, at, direction))}
            onDelete={() => actions.remove.mutate(room.id)}
          />
          <span className="text-caption text-muted tabular">
            {inRoom.length === 0 ? "empty" : `${on} of ${inRoom.length} on`}
          </span>
        </div>

        {switchable.length >= 2 ? (
          <div className="flex shrink-0 gap-1.5">
            <Button
              variant="ghost"
              size="sm"
              disabled={on === 0}
              onClick={() => void setRoomPower(inRoom, false)}
            >
              All off
            </Button>
            <Button variant="secondary" size="sm" onClick={() => void setRoomPower(inRoom, true)}>
              All on
            </Button>
          </div>
        ) : null}
      </header>

      {inRoom.length === 0 ? (
        <p className="text-small text-muted">No devices in here yet.</p>
      ) : (
        <div className="grid grid-cols-[repeat(auto-fill,minmax(210px,1fr))] gap-3">
          {inRoom.map((device) => (
            <DeviceCard
              key={device.id}
              device={device}
              disabled={device.status === "offline"}
              onToggle={() => control(device, TRAIT_POWER, !isOn(device))}
            />
          ))}
        </div>
      )}
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

function moved<T>(list: T[], index: number, direction: -1 | 1): T[] {
  const target = index + direction;
  if (target < 0 || target >= list.length) return list;

  const next = [...list];
  const [lifted] = next.splice(index, 1);
  if (lifted !== undefined) next.splice(target, 0, lifted);

  return next;
}
