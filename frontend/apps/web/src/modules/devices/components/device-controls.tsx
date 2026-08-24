import { Button, Slider, Switch, cn } from "@ava/ui";
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
 * `levelOverride` reports the in-flight slider value so the caller can light the
 * device on the stage while a drag is still happening.
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
    Math.min(
      Math.max(numberOf(device, TRAIT_BRIGHTNESS) ?? dimming.max, dimming.min),
      dimming.max,
    ),
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
    <div className="grid gap-5">
      <div className="flex items-center justify-between">
        <Heading>Power</Heading>
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
        <div className="grid gap-2.5">
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

      <div className="grid gap-4 border-t border-border pt-4">
        <Button
          variant="ghost"
          size="sm"
          aria-expanded={details}
          onClick={() => setDetails((open) => !open)}
          className="justify-between px-0 hover:bg-transparent"
        >
          Details
          <ChevronDownIcon
            aria-hidden
            className={cn(
              "size-4 transition-transform duration-200 ease-out",
              details && "rotate-180",
            )}
          />
        </Button>

        {details ? (
          <div className="grid gap-5">
            <RoomPicker device={device} />
            {device.kind === "plug" ? <AppliancePicker device={device} /> : null}
            <div className="grid gap-1 text-caption text-subtle">
              <span>{profile.form}</span>
              <span className="font-mono">{device.external_id}</span>
              {device.vendor ? <span>{device.vendor}</span> : null}
            </div>
          </div>
        ) : null}
      </div>
    </div>
  );
}
