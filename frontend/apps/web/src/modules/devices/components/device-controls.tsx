import { Slider, Switch, cn } from "@ava/ui";
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
import { ChevronDownIcon } from "lucide-react";
import { useState } from "react";

import { AppliancePicker } from "./appliance-picker";
import { ColorControl } from "./color-control";
import { RoomPicker } from "./room-picker";
import { TraitControl, TraitReading } from "./trait-control";
import { deviceColor } from "./device-stage";
import { useDeviceControl } from "../use-devices";
import { useLiveSlider } from "../use-live-slider";

function Heading({ children }: { children: string }) {
  return (
    <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
      {children}
    </span>
  );
}

/** A fact about the device rather than something you can change. */
function Fact({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex items-baseline justify-between gap-4">
      <span className="shrink-0 text-caption text-subtle">{label}</span>
      <span className={cn("min-w-0 truncate text-small text-muted", mono && "font-mono")}>
        {value}
      </span>
    </div>
  );
}

/**
 * Everything you touch on a device, and nothing you don't.
 *
 * The old device page put power, brightness and colour in the same column as
 * the room assignment, the appliance type and the vendor id. Those are used at
 * completely different rates — several times a day against once, ever — and
 * mixing them is what made the panel too tall to fit a phone sheet without
 * scrolling inside something that also drags to dismiss. Configuration is
 * folded away behind Details.
 *
 * `onLevelChange` reports the in-flight slider value so the caller can light
 * the device on the stage while a drag is still happening.
 */
export function DeviceControls({
  device,
  offline,
  onLevelChange,
}: {
  device: DeviceDto;
  offline: boolean;
  onLevelChange?: (level: number | null) => void;
}) {
  const control = useDeviceControl();
  const [details, setDetails] = useState(false);

  const dimming = brightnessRange(device) ?? { min: 0, max: 100, step: 1, unit: "%" };
  const warmth = kelvinRange(device);

  const send = (trait: string, value: TraitValue) => control(device, trait, value);

  const brightness = useLiveSlider(
    Math.min(Math.max(numberOf(device, TRAIT_BRIGHTNESS) ?? dimming.max, dimming.min), dimming.max),
    (value) => {
      onLevelChange?.(value);
      send(TRAIT_BRIGHTNESS, value);
    },
    (value) => {
      onLevelChange?.(null);
      send(TRAIT_BRIGHTNESS, value);
    },
  );

  const on = isOn(device);
  const dimmable = supports(device, TRAIT_BRIGHTNESS);
  const extras = writableTraits(device).filter(
    (capability) =>
      capability.trait !== TRAIT_POWER &&
      capability.trait !== TRAIT_BRIGHTNESS &&
      capability.trait !== TRAIT_COLOR_TEMP,
  );
  const sensors = readings(device);
  const profile = deviceProfile(device);

  return (
    <div className="grid gap-6">
      <div className="flex min-h-11 items-center justify-between gap-4">
        <Heading>Power</Heading>
        <Switch
          checked={on}
          disabled={offline}
          onCheckedChange={(next) => send(TRAIT_POWER, next)}
          aria-label={`${device.name} power`}
        />
      </div>

      {dimmable ? (
        <div className="grid gap-3">
          <div className="flex items-baseline justify-between gap-4">
            <Heading>Brightness</Heading>
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
        <div className="grid gap-3">
          <Heading>Colour</Heading>
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
        <div className="grid gap-3 border-t border-border pt-6">
          {sensors.map((capability) => (
            <TraitReading
              key={capability.trait}
              capability={capability}
              value={device.state[capability.trait]}
            />
          ))}
        </div>
      ) : null}

      <div className="border-t border-border pt-2">
        <button
          type="button"
          aria-expanded={details}
          onClick={() => setDetails((open) => !open)}
          className={cn(
            "-mx-2 flex min-h-11 w-[calc(100%+1rem)] items-center justify-between gap-4 rounded-sm px-2",
            "text-small font-semibold text-muted",
            "transition-colors duration-150 ease-out hover:text-fg",
            "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
          )}
        >
          Details
          <ChevronDownIcon
            aria-hidden
            className={cn(
              "size-4 shrink-0 transition-transform duration-200 ease-out",
              details && "rotate-180",
            )}
          />
        </button>

        {/* Animating grid-template-rows from 0fr to 1fr is the one way to slide
            a panel of unknown height open without measuring it in JavaScript. */}
        <div
          className={cn(
            "grid transition-[grid-template-rows] duration-300 ease-out-soft",
            "motion-reduce:transition-none",
            details ? "grid-rows-[1fr]" : "grid-rows-[0fr]",
          )}
        >
          <div className="overflow-hidden">
            <div className="grid gap-5 pb-1 pt-4">
              <RoomPicker device={device} />
              {device.kind === "plug" ? <AppliancePicker device={device} /> : null}

              <div className="grid gap-2.5 rounded-md border border-border bg-surface p-3.5">
                <Fact label="Type" value={profile.form} />
                {device.vendor ? <Fact label="Vendor" value={device.vendor} /> : null}
                <Fact label="ID" value={device.external_id} mono />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
