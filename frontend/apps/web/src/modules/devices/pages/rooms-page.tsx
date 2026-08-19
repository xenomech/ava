import { Device, DeviceHalo, Switch, cn } from "@ava/ui";
import { Link } from "@tanstack/react-router";

import { useDeviceStore, useDevices } from "../store";

export function RoomsPage() {
  const { devices, focus } = useDevices();
  const toggle = useDeviceStore((s) => s.toggle);

  const rooms = [...new Set(devices.map((d) => d.room))];

  return (
    <div className="grid gap-8 p-5 sm:p-6">
      <header>
        <h1 className="text-display font-semibold">Rooms</h1>
        <p className="mt-1 text-small text-muted">
          {rooms.length} rooms · {devices.length} devices
        </p>
      </header>

      {rooms.map((room) => {
        const inRoom = devices.filter((d) => d.room === room);
        const on = inRoom.filter((d) => d.level > 0).length;

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
                <article
                  key={device.id}
                  className={cn(
                    "grid gap-3 rounded-lg border border-border bg-surface p-4",
                    "transition-colors duration-150 ease-out",
                    device.level > 0 && "border-border-strong",
                  )}
                  style={{ "--level": device.level, "--lit": device.color } as React.CSSProperties}
                >
                  <Link
                    to="/"
                    onClick={() => focus(device.id)}
                    className="relative grid h-32 place-items-center rounded-sm"
                    aria-label={`Open ${device.name}`}
                  >
                    <DeviceHalo className="w-3/4" />
                    <Device
                      kind={device.kind}
                      level={device.level}
                      color={device.color}
                      className="h-[86%]"
                    />
                  </Link>

                  <div className="flex items-end justify-between gap-3">
                    <span className="min-w-0">
                      <span className="block truncate text-body font-semibold">{device.name}</span>
                      <span className="block font-mono text-caption text-subtle tabular">
                        {device.level > 0 ? `${device.level}%` : "Off"}
                      </span>
                    </span>

                    <Switch
                      checked={device.level > 0}
                      onCheckedChange={() => toggle(device.id)}
                      aria-label={device.name}
                    />
                  </div>
                </article>
              ))}
            </div>
          </section>
        );
      })}
    </div>
  );
}
