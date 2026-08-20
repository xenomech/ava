import { Button, Chip, Slider, Switch, cn } from "@ava/ui";
import type { DeviceAction, DeviceDto } from "@ava/contracts";
import { useQuery } from "@tanstack/react-query";
import { Device } from "@ava/ui";
import { useNavigate, useSearch } from "@tanstack/react-router";

import { hubQueries } from "@/modules/hub";
import { sendCommand } from "../api";
import { Loader } from "@/shared/components/loader";
import { ColorControl } from "../components/color-control";
import { DeviceStage, deviceColor, deviceKind, deviceLevel } from "../components/device-stage";
import { NoDevices } from "../components/empty-state";
import { HubOfflineNotice } from "../components/hub-notice";
import { useDeviceCommand, useDevices } from "../use-devices";
import { useLiveSlider } from "../use-live-slider";

export function ConsolePage() {
  const { devices, isPending } = useDevices();
  const hubs = useQuery(hubQueries.list());
  const command = useDeviceCommand();

  const navigate = useNavigate();
  const search = useSearch({ strict: false }) as { device?: string };

  const device = devices.find((entry) => entry.id === search.device) ?? devices[0];

  const limits = device?.state.limits;
  const brightnessMin = limits?.brightness_min ?? 0;
  const brightnessMax = limits?.brightness_max ?? 100;
  const kelvinMin = limits?.kelvin_min;
  const kelvinMax = limits?.kelvin_max;

  const send = (action: DeviceAction, value: number) => {
    if (!device) return;

    void sendCommand(device.id, { action, value }).catch(() => undefined);
  };

  const settle = (action: DeviceAction, value: number) => {
    if (!device) return;

    command.mutate({ device, action, value });
  };

  const brightness = useLiveSlider(
    clamp(device?.state.brightness ?? brightnessMax, brightnessMin, brightnessMax),
    (value) => send("brightness", value),
    (value) => settle("brightness", value),
  );

  const focus = (id: string) => {
    void navigate({ to: "/", search: { device: id }, replace: true });
  };

  if (isPending) return <Loader label="Loading devices" />;

  if (!device) return <NoDevices hasHub={(hubs.data ?? []).length > 0} />;

  const hub = (hubs.data ?? []).find((entry) => entry.id === device.hub_id);
  const hubOffline = hub !== undefined && !hub.online;
  const offline = device.status === "offline" || hubOffline;
  const capabilities = device.state.capabilities;
  const canWarm =
    capabilities.includes("color_temp") && kelvinMin !== undefined && kelvinMax !== undefined;

  return (
    <div
      className={cn(
        "grid h-full",
        hubOffline ? "grid-rows-[auto_minmax(0,1fr)_auto]" : "grid-rows-[minmax(0,1fr)_auto]",
      )}
    >
      {hubOffline ? <HubOfflineNotice name={hub.name} /> : null}

      <div className="min-h-0 overflow-y-auto lg:grid lg:grid-cols-[minmax(0,1fr)_330px] lg:overflow-hidden">
        <main className="grid min-h-0 grid-rows-[34vh_auto] p-5 sm:p-6 lg:grid-rows-[minmax(0,1fr)_auto]">
          <DeviceStage
            devices={devices}
            focusedID={device.id}
            levelOverride={brightness.dragging}
          />

          <div className="flex items-end justify-between gap-4 pt-4">
            <div>
              <h1 className="text-display font-semibold">{device.name}</h1>
              <p className="mt-1.5 flex items-center gap-2 text-small text-muted">
                {device.room || "No room"}
                {offline ? (
                  <Chip tone="warning">{hubOffline ? "Hub offline" : "Offline"}</Chip>
                ) : null}
              </p>
            </div>

            <div className="flex items-start gap-0.5" data-slot="reading">
              {device.state.power || brightness.dragging !== null ? (
                <>
                  <b className="text-hero font-semibold">{brightness.value}</b>
                  <span className="mt-1 text-small text-subtle">%</span>
                </>
              ) : (
                <b className="text-hero font-semibold text-subtle">Off</b>
              )}
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
                <output className="font-mono text-small tabular">{brightness.value}%</output>
              </div>
              <Slider
                value={[brightness.value]}
                min={brightnessMin}
                max={brightnessMax}
                step={1}
                lit
                disabled={offline}
                className={cn(!device.state.power && "opacity-40")}
                aria-label="Brightness"
                onValueChange={([value]) => brightness.change(value ?? brightnessMin)}
                onValueCommit={([value]) => brightness.release(value ?? brightnessMin)}
                style={{ "--lit": deviceColor(device) } as React.CSSProperties}
              />
            </div>
          ) : null}

          {canWarm ? (
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
                disabled={offline}
                onWhitePreview={(kelvin) => send("color_temp", kelvin)}
                onWhite={(kelvin) => settle("color_temp", kelvin)}
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

function clamp(value: number, low: number, high: number) {
  return Math.min(Math.max(value, low), high);
}
