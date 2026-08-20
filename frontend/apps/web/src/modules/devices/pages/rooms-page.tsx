import { Chip, Device, DeviceHalo, Switch, cn } from "@ava/ui";
import type { DeviceDto } from "@ava/contracts";
import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";

import { hubQueries } from "@/modules/hub";
import { Loader } from "@/shared/components/loader";
import { deviceColor, deviceKind, deviceLevel } from "../components/device-stage";
import { NoDevices } from "../components/empty-state";
import { useDeviceCommand, useDevices } from "../use-devices";

const UNASSIGNED = "No room";

export function RoomsPage() {
  const { devices, isPending } = useDevices();
  const hubs = useQuery(hubQueries.list());
  const command = useDeviceCommand();

  if (isPending) return <Loader label="Loading devices" />;

  if (devices.length === 0) return <NoDevices hasHub={(hubs.data ?? []).length > 0} />;

  const rooms = [...new Set(devices.map((device) => device.room || UNASSIGNED))];

  return (
    <div className="grid gap-8 p-5 sm:p-6">
      <header>
        <h1 className="text-display font-semibold">Rooms</h1>
        <p className="mt-1 text-small text-muted">
          {rooms.length} {rooms.length === 1 ? "room" : "rooms"} · {devices.length}{" "}
          {devices.length === 1 ? "device" : "devices"}
        </p>
      </header>

      {rooms.map((room) => {
        const inRoom = devices.filter((device) => (device.room || UNASSIGNED) === room);
        const on = inRoom.filter((device) => device.state.power).length;

        return (
          <section key={room} className="grid gap-3">
            <div className="flex items-baseline justify-between">
              <h2 className="text-title font-semibold">{room}</h2>
              <span className="font-mono text-caption text-subtle tabular">
                {on} of {inRoom.length} on
              </span>
            </div>

            <div className="grid grid-cols-[repeat(auto-fill,minmax(210px,1fr))] gap-3">
              {inRoom.map((device) => (
                <DeviceCard
                  key={device.id}
                  device={device}
                  disabled={device.status === "offline" || command.isPending}
                  onToggle={() =>
                    command.mutate({ device, action: "power", value: !device.state.power })
                  }
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
        device.state.power && "border-border-strong",
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
            checked={device.state.power}
            disabled={disabled}
            onCheckedChange={onToggle}
            aria-label={device.name}
          />
        )}
      </div>
    </article>
  );
}
