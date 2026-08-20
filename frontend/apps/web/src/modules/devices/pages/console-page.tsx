import { Button, Chip, Slider, Switch, cn } from "@ava/ui";
import type { DeviceDto } from "@ava/contracts";
import { useQuery } from "@tanstack/react-query";
import { Device } from "@ava/ui";
import { useNavigate, useSearch } from "@tanstack/react-router";
import { useCallback, useState } from "react";

import { hubQueries } from "@/modules/hub";
import { sendCommand } from "../api";
import { Loader } from "@/shared/components/loader";
import { ColorControl } from "../components/color-control";
import { DeviceStage, deviceColor, deviceKind, deviceLevel } from "../components/device-stage";
import { NoDevices } from "../components/empty-state";
import { useDeviceCommand, useDevices } from "../use-devices";
import { useThrottled } from "../use-throttled-command";

export function ConsolePage() {
  const { devices, isPending } = useDevices();
  const hubs = useQuery(hubQueries.list());
  const command = useDeviceCommand();

  const navigate = useNavigate();
  const [dragging, setDragging] = useState<number | null>(null);

  const sendLive = useThrottled<{ id: string; value: number }>(
    useCallback(({ id, value }) => {
      void sendCommand(id, { action: "brightness", value }).catch(() => undefined);
    }, []),
  );
  const search = useSearch({ strict: false }) as { device?: string };

  const focus = (id: string) => {
    void navigate({ to: "/", search: { device: id }, replace: true });
  };

  if (isPending) return <Loader label="Loading devices" />;

  if (devices.length === 0) return <NoDevices hasHub={(hubs.data ?? []).length > 0} />;

  const device = devices.find((entry) => entry.id === search.device) ?? devices[0]!;
  const level = dragging ?? deviceLevel(device);
  const offline = device.status === "offline";
  const capabilities = device.state.capabilities;
  const limits = device.state.limits;
  const brightnessMin = limits?.brightness_min ?? 0;
  const brightnessMax = limits?.brightness_max ?? 100;
  const kelvinMin = limits?.kelvin_min;
  const kelvinMax = limits?.kelvin_max;

  return (
    <div className="grid h-full grid-rows-[minmax(0,1fr)_auto]">
      <div className="min-h-0 overflow-y-auto lg:grid lg:grid-cols-[minmax(0,1fr)_330px] lg:overflow-hidden">
        <main className="grid min-h-0 grid-rows-[34vh_auto] p-5 sm:p-6 lg:grid-rows-[minmax(0,1fr)_auto]">
          <DeviceStage devices={devices} focusedID={device.id} levelOverride={dragging} />

          <div className="flex items-end justify-between gap-4 pt-4">
            <div>
              <h1 className="text-display font-semibold">{device.name}</h1>
              <p className="mt-1.5 flex items-center gap-2 text-small text-muted">
                {device.room || "No room"}
                {offline ? <Chip tone="warning">Offline</Chip> : null}
              </p>
            </div>

            <div className="flex items-start gap-0.5" data-slot="reading">
              <b className="text-hero font-semibold">{level}</b>
              <span className="mt-1 text-small text-subtle">%</span>
            </div>
          </div>
        </main>

        <aside className="flex flex-col gap-5 border-t border-border p-5 lg:min-h-0 lg:overflow-y-auto lg:border-t-0 lg:border-l">
          <div className="flex items-center justify-between">
            <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
              Power
            </span>
            <Switch
              checked={device.state.power}
              disabled={offline || command.isPending}
              onCheckedChange={(on) => command.mutate({ device, action: "power", value: on })}
              aria-label={`${device.name} power`}
            />
          </div>

          {capabilities.includes("brightness") ? (
            <div className="grid gap-2.5">
              <div className="flex items-baseline justify-between">
                <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
                  Brightness
                </span>
                <output className="font-mono text-small tabular">{level}%</output>
              </div>
              <Slider
                value={[level]}
                min={brightnessMin}
                max={brightnessMax}
                step={1}
                lit
                disabled={offline}
                aria-label="Brightness"
                onValueChange={([value]) => {
                  setDragging(value ?? 0);
                  sendLive({ id: device.id, value: value ?? 0 });
                }}
                onValueCommit={([value]) => {
                  setDragging(null);
                  command.mutate({ device, action: "brightness", value: value ?? 0 });
                }}
                style={{ "--lit": deviceColor(device) } as React.CSSProperties}
              />
            </div>
          ) : null}

          {capabilities.includes("color_temp") ? (
            <div className="grid gap-2.5">
              <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
                Colour
              </span>
              <ColorControl
                color={deviceColor(device)}
                kelvin={device.state.color_temp ?? null}
                kelvinMin={kelvinMin}
                kelvinMax={kelvinMax}
                showColor={capabilities.includes("color")}
                onWhite={(kelvin) =>
                  command.mutate({ device, action: "color_temp", value: kelvin })
                }
                onColor={() => undefined}
              />
            </div>
          ) : null}

          <div className="grid gap-1 text-caption text-subtle">
            <span className="font-mono">{device.external_id}</span>
            {device.state.vendor ? <span>{device.state.vendor}</span> : null}
          </div>
        </aside>
      </div>

      <footer
        className="flex gap-2 overflow-x-auto border-t border-border p-3"
        aria-label="All devices"
      >
        {devices.map((entry) => (
          <DeviceTile
            key={entry.id}
            device={entry}
            focused={entry.id === device.id}
            onFocus={() => focus(entry.id)}
          />
        ))}
      </footer>
    </div>
  );
}

function DeviceTile({
  device,
  focused,
  onFocus,
}: {
  device: DeviceDto;
  focused: boolean;
  onFocus: () => void;
}) {
  const level = deviceLevel(device);

  return (
    <Button
      variant="ghost"
      aria-current={focused}
      onClick={onFocus}
      className={cn(
        "h-auto w-32 shrink-0 flex-col items-start gap-1.5 rounded-md border-border p-2.5",
        focused && "border-fg",
        device.status === "offline" && "opacity-50",
      )}
      style={{ "--level": level, "--lit": deviceColor(device) } as React.CSSProperties}
    >
      <span className="grid h-14 w-full place-items-center">
        <Device
          kind={deviceKind(device)}
          level={level}
          color={deviceColor(device)}
          className="h-full"
        />
      </span>
      <span className="w-full truncate text-left text-small font-semibold">{device.name}</span>
      <span className="font-mono text-caption font-normal text-subtle tabular">
        {device.status === "offline" ? "Offline" : level > 0 ? `${level}%` : "Off"}
      </span>
    </Button>
  );
}
