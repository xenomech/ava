import { Button, Chip, Slider, Switch, cn } from "@ava/ui";
import {
  TRAIT_BRIGHTNESS,
  TRAIT_COLOR,
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  brightnessRange,
  deviceProfile,
  hasColor,
  isOn,
  kelvinRange,
  numberOf,
  readings,
  supports,
  writableTraits,
  type DeviceDto,
  type TraitValue,
} from "@ava/contracts";
import { useQuery } from "@tanstack/react-query";
import { Device } from "@ava/ui";
import { useNavigate, useParams } from "@tanstack/react-router";

import { hubQueries } from "@/modules/hub";
import { Loader } from "@/shared/components/loader";
import { AppliancePicker } from "../components/appliance-picker";
import { RoomPicker } from "../components/room-picker";
import { ColorControl } from "../components/color-control";
import {
  DeviceStage,
  OnAPlug,
  deviceColor,
  deviceKind,
  deviceLabel,
  deviceLevel,
} from "../components/device-stage";
import { Missing, NoDevices } from "../components/empty-state";
import { TraitControl, TraitReading } from "../components/trait-control";
import { HubOfflineNotice } from "../components/hub-notice";
import { useDeviceControl, useDevices } from "../use-devices";
import { useLiveSlider } from "../use-live-slider";

/** One device, reached from a room's strip or from the tray along the foot. */
export function DevicePage() {
  const { deviceId } = useParams({ from: "/_protected/devices/$deviceId" });
  const { devices, isPending } = useDevices();
  const hubs = useQuery(hubQueries.list());
  const control = useDeviceControl();

  const navigate = useNavigate();

  const device = devices.find((entry) => entry.id === deviceId);

  const dimming = brightnessRange(device) ?? { min: 0, max: 100, step: 1, unit: "%" };
  const warmth = kelvinRange(device);

  const send = (trait: string, value: TraitValue) => {
    if (!device) return;

    control(device, trait, value);
  };

  const brightness = useLiveSlider(
    clamp(numberOf(device, TRAIT_BRIGHTNESS) ?? dimming.max, dimming.min, dimming.max),
    (value) => send(TRAIT_BRIGHTNESS, value),
    (value) => send(TRAIT_BRIGHTNESS, value),
  );

  const focus = (id: string) => {
    void navigate({ to: "/devices/$deviceId", params: { deviceId: id }, replace: true });
  };

  if (isPending) return <Loader label="Loading device" />;

  if (!device) {
    return devices.length === 0 ? (
      <NoDevices hasHub={(hubs.data ?? []).length > 0} />
    ) : (
      <Missing
        title="That device is gone"
        detail="It may have been removed, or it belongs to a hub you no longer have."
      />
    );
  }

  const on = isOn(device);
  const dimmable = supports(device, TRAIT_BRIGHTNESS);
  const extras = writableTraits(device).filter(
    (capability) =>
      capability.trait !== TRAIT_POWER &&
      capability.trait !== TRAIT_BRIGHTNESS &&
      capability.trait !== TRAIT_COLOR_TEMP,
  );
  const sensors = readings(device);

  const hub = (hubs.data ?? []).find((entry) => entry.id === device.hub_id);
  const hubOffline = hub !== undefined && !hub.online;
  const offline = device.status === "offline" || hubOffline;
  const profile = deviceProfile(device);

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
              {!on ? (
                <b className="text-hero font-semibold text-subtle">Off</b>
              ) : dimmable ? (
                <>
                  <b className="text-hero font-semibold">{brightness.value}</b>
                  <span className="mt-1 text-small text-subtle">%</span>
                </>
              ) : (
                <b className="text-hero font-semibold">On</b>
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
              checked={on}
              disabled={offline}
              onCheckedChange={(next) => send(TRAIT_POWER, next)}
              aria-label={`${device.name} power`}
            />
          </div>

          {dimmable ? (
            <div className="grid gap-2.5">
              <div className="flex items-baseline justify-between">
                <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
                  Brightness
                </span>
                <output className="font-mono text-small tabular">{brightness.value}%</output>
              </div>
              <Slider
                value={[brightness.value]}
                min={dimming.min}
                max={dimming.max}
                step={dimming.step}
                lit
                disabled={offline}
                className={cn(!on && "opacity-40")}
                aria-label="Brightness"
                onValueChange={([value]) => brightness.change(value ?? dimming.min)}
                onValueCommit={([value]) => brightness.release(value ?? dimming.min)}
                style={{ "--lit": deviceColor(device) } as React.CSSProperties}
              />
            </div>
          ) : null}

          {warmth ? (
            <div className="grid gap-2.5">
              <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
                Colour
              </span>
              <ColorControl
                color={deviceColor(device)}
                kelvin={numberOf(device, TRAIT_COLOR_TEMP) ?? null}
                kelvinMin={warmth.min}
                kelvinMax={warmth.max}
                showColor={hasColor(device)}
                disabled={offline}
                onWhitePreview={(kelvin) => send(TRAIT_COLOR_TEMP, kelvin)}
                onWhite={(kelvin) => send(TRAIT_COLOR_TEMP, kelvin)}
                onColor={(hex) => send(TRAIT_COLOR, hex)}
              />
            </div>
          ) : null}

          {extras.map((capability) => (
            <TraitControl
              key={capability.trait}
              capability={capability}
              value={device.state[capability.trait]}
              disabled={offline}
              onChange={(value) => send(capability.trait, value)}
            />
          ))}

          {sensors.length > 0 ? (
            <div className="grid gap-2 border-t border-border pt-4">
              {sensors.map((capability) => (
                <TraitReading
                  key={capability.trait}
                  capability={capability}
                  value={device.state[capability.trait]}
                />
              ))}
            </div>
          ) : null}

          <RoomPicker device={device} />

          {device.kind === "plug" ? <AppliancePicker device={device} /> : null}

          <div className="grid gap-1 text-caption text-subtle">
            <span>{profile.form}</span>
            <span className="font-mono">{device.external_id}</span>
            {device.vendor ? <span>{device.vendor}</span> : null}
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
        "relative h-auto w-32 shrink-0 flex-col items-start gap-1.5 rounded-md border-border p-2.5",
        focused && "border-fg",
        device.status === "offline" && "opacity-50",
      )}
      style={{ "--level": level, "--lit": deviceColor(device) } as React.CSSProperties}
    >
      {device.kind === "plug" && device.appliance ? (
        <OnAPlug className="absolute right-2 top-2 size-5" />
      ) : null}
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
        {deviceLabel(device)}
      </span>
    </Button>
  );
}

function clamp(value: number, low: number, high: number) {
  return Math.min(Math.max(value, low), high);
}
