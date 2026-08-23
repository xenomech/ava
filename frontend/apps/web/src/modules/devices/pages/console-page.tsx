import {
  Button,
  Chip,
  Drawer,
  DrawerBody,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
  cn,
} from "@ava/ui";
import {
  TRAIT_BRIGHTNESS,
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  brightnessRange,
  isOn,
  numberOf,
  supports,
  type TraitValue,
} from "@ava/contracts";
import { useQuery } from "@tanstack/react-query";
import { useNavigate, useSearch } from "@tanstack/react-router";
import { useState, type CSSProperties } from "react";

import { hubQueries } from "@/modules/hub";
import { Loader } from "@/shared/components/loader";
import { FixtureStrip, RoomRail, RoomTabs } from "../components/console-rails";
import { DeviceControls } from "../components/device-controls";
import { DeviceHero } from "../components/device-hero";
import { deviceColor, deviceLevel } from "../device-view";
import { NoDevices } from "../components/empty-state";
import { HubOfflineNotice } from "../components/hub-notice";
import { useDeviceControl, useDevices, useRoomPower } from "../use-devices";
import { useLiveSlider } from "../use-live-slider";
import { useRoomGroups } from "../use-room-groups";

export function ConsolePage() {
  const { devices, isPending } = useDevices();
  const hubs = useQuery(hubQueries.list());
  const groups = useRoomGroups(devices);
  const control = useDeviceControl();
  const setRoomPower = useRoomPower();

  const navigate = useNavigate();
  const search = useSearch({ strict: false }) as { device?: string };
  const [sheetOpen, setSheetOpen] = useState(false);

  // The room is derived from the focused fixture rather than tracked separately,
  // so the two can never disagree about where you are.
  const focused =
    groups.flatMap((group) => group.devices).find((entry) => entry.id === search.device) ??
    groups[0]?.devices[0];

  const group =
    groups.find((entry) => entry.devices.some((d) => d.id === focused?.id)) ?? groups[0];

  const dimming = brightnessRange(focused) ?? { min: 0, max: 100, step: 1, unit: "%" };

  const send = (trait: string, value: TraitValue) => {
    if (!focused) return;

    control(focused, trait, value);
  };

  const brightness = useLiveSlider(
    focused ? clamp(numberOf(focused, TRAIT_BRIGHTNESS) ?? dimming.max, dimming) : dimming.max,
    (value) => send(TRAIT_BRIGHTNESS, value),
    (value) => send(TRAIT_BRIGHTNESS, value),
  );

  if (isPending) return <Loader label="Loading devices" />;

  if (!focused || !group) return <NoDevices hasHub={(hubs.data ?? []).length > 0} />;

  const focus = (id: string) => void navigate({ to: "/", search: { device: id }, replace: true });

  // Flicking wraps inside the room. Leaving it is a deliberate move on the rail,
  // so you cannot drift into another room without noticing.
  const step = (direction: 1 | -1) => {
    const list = group.devices;
    const at = list.findIndex((entry) => entry.id === focused.id);
    const next = list[(at + direction + list.length) % list.length];

    if (next) focus(next.id);
  };

  const pickRoom = (key: string) => {
    const first = groups.find((entry) => entry.key === key)?.devices[0];

    if (first) focus(first.id);
  };

  const hub = (hubs.data ?? []).find((entry) => entry.id === focused.hub_id);
  const hubOffline = hub !== undefined && !hub.online;
  const offline = focused.status === "offline" || hubOffline;

  const on = isOn(focused);
  const dimmable = supports(focused, TRAIT_BRIGHTNESS);
  const level = brightness.dragging ?? deviceLevel(focused);
  const kelvin = numberOf(focused, TRAIT_COLOR_TEMP);

  const controls = (
    <DeviceControls
      device={focused}
      brightness={brightness.value}
      dimming={dimming}
      dimmable={dimmable}
      offline={offline}
      onSend={send}
      onBrightnessChange={brightness.change}
      onBrightnessCommit={brightness.release}
    />
  );

  const fixtures = <FixtureStrip devices={group.devices} focusedID={focused.id} onFocus={focus} />;

  return (
    // The notice sits outside the console grid on purpose: as a grid child it
    // competed for a cell with the rail and the stage and drew over both.
    <div className="flex h-full flex-col">
      {hubOffline ? <HubOfflineNotice name={hub.name} /> : null}

      <div className="grid min-h-0 flex-1 lg:grid-cols-[204px_minmax(0,1fr)]">
        <aside className="hidden min-h-0 overflow-y-auto border-r border-border p-3 lg:block">
          <RoomRail groups={groups} activeKey={group.key} onPick={pickRoom} />
        </aside>

        <div className="grid min-h-0 grid-rows-[auto_minmax(0,1fr)_auto] xl:grid-cols-[minmax(0,1fr)_344px]">
          {/* On a phone the tabs are how you see which room you are in, so the
            heading below is desktop-only rather than saying it twice. */}
          <RoomTabs
            groups={groups}
            activeKey={group.key}
            onPick={pickRoom}
            className="px-5 pt-4 pb-1 sm:px-6 lg:hidden"
          />

          <header className="col-start-1 flex items-center justify-between gap-4 px-5 pt-3 sm:px-6 lg:pt-4">
            <div className="min-w-0">
              <h1 className="hidden truncate text-title font-semibold lg:block">{group.name}</h1>
              <p className="text-caption text-muted tabular lg:mt-0.5">
                {group.on} of {group.devices.length} on
              </p>
            </div>

            <div className="flex shrink-0 gap-1.5">
              <Button
                variant="ghost"
                size="sm"
                disabled={group.on === 0}
                onClick={() => void setRoomPower(group.devices, false)}
              >
                All off
              </Button>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => void setRoomPower(group.devices, true)}
              >
                All on
              </Button>
            </div>
          </header>

          <main className="col-start-1 grid min-h-0 grid-rows-[minmax(0,1fr)_auto] gap-4 px-5 py-4 sm:px-6">
            <DeviceHero
              device={focused}
              level={level}
              dimmable={dimmable}
              offline={offline}
              onToggle={() => send(TRAIT_POWER, !on)}
              onDim={brightness.change}
              onDimEnd={brightness.release}
              onStep={step}
              onOpenSheet={() => setSheetOpen(true)}
            />

            <div className="flex items-end justify-between gap-4">
              <div className="min-w-0">
                <h2 className="truncate text-display font-semibold">{focused.name}</h2>
                {offline ? (
                  <Chip tone="warning" className="mt-1.5">
                    {hubOffline ? "Hub offline" : "Offline"}
                  </Chip>
                ) : null}
              </div>

              <p className="flex shrink-0 items-start gap-0.5" data-slot="reading">
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
              </p>
            </div>
          </main>

          <footer className="col-start-1 grid gap-3 px-5 pb-4 sm:px-6 xl:hidden">
            {fixtures}

            <Drawer open={sheetOpen} onOpenChange={setSheetOpen}>
              <DrawerTrigger asChild>
                <button
                  type="button"
                  aria-label={`Adjust ${focused.name}`}
                  style={{ "--lit": deviceColor(focused) } as CSSProperties}
                  className={cn(
                    "flex items-center justify-between gap-3 rounded-lg border border-border px-4 py-3",
                    "transition-colors duration-150 ease-out hover:border-border-strong",
                  )}
                >
                  <span className="flex items-center gap-2.5">
                    <span aria-hidden className="h-1 w-8 rounded-full bg-border-strong" />
                    <span className="font-mono text-small text-muted tabular">
                      {kelvin ? `${kelvin}K · ` : ""}
                      {dimmable ? `${brightness.value}%` : on ? "On" : "Off"}
                    </span>
                  </span>
                  <span className="text-small font-semibold">Adjust</span>
                </button>
              </DrawerTrigger>

              <DrawerContent>
                <DrawerHeader>
                  <DrawerTitle>{focused.name}</DrawerTitle>
                </DrawerHeader>
                <DrawerBody className="pb-8">{controls}</DrawerBody>
              </DrawerContent>
            </Drawer>
          </footer>

          <aside className="hidden min-h-0 overflow-y-auto border-l border-border p-5 xl:col-start-2 xl:row-span-3 xl:row-start-1 xl:block">
            <div className="grid gap-5">
              <FixtureStrip
                devices={group.devices}
                focusedID={focused.id}
                onFocus={focus}
                layout="grid"
              />
              {controls}
              <KeyboardMap />
            </div>
          </aside>
        </div>
      </div>
    </div>
  );
}

const KEYS = [
  ["Space", "Power"],
  ["↑ ↓", "Brightness"],
  ["← →", "Fixture"],
] as const;

function KeyboardMap() {
  return (
    <div className="grid gap-2 border-t border-border pt-4 pb-1">
      {KEYS.map(([key, label]) => (
        <div key={key} className="flex items-center justify-between gap-3">
          <kbd className="rounded-xs border border-border px-1.5 py-0.5 font-mono text-caption text-subtle">
            {key}
          </kbd>
          <span className="text-caption text-subtle">{label}</span>
        </div>
      ))}
    </div>
  );
}

function clamp(value: number, range: { min: number; max: number }) {
  return Math.min(Math.max(value, range.min), range.max);
}
