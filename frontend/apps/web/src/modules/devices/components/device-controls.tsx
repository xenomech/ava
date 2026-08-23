import { Slider, Switch, cn } from "@ava/ui";
import {
  TRAIT_BRIGHTNESS,
  TRAIT_COLOR,
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  deviceProfile,
  hasColor,
  isOn,
  kelvinRange,
  numberOf,
  readings,
  writableTraits,
  type DeviceDto,
  type Range,
  type TraitValue,
} from "@ava/contracts";

import { AppliancePicker } from "./appliance-picker";
import { ColorControl } from "./color-control";
import { RenameDevice } from "./rename-device";
import { RoomPicker } from "./room-picker";
import { TraitControl, TraitReading } from "./trait-control";
import { deviceColor } from "../device-view";

// The full control surface for one fixture. On a phone it is the contents of
// the pull-up sheet; on a desktop the same component is docked open. Warmth
// sits below brightness on purpose — it is set once per fixture, where power
// and brightness are daily.
export function DeviceControls({
  device,
  brightness,
  dimming,
  dimmable,
  offline,
  onSend,
  onBrightnessChange,
  onBrightnessCommit,
}: {
  device: DeviceDto;
  brightness: number;
  dimming: Range;
  dimmable: boolean;
  offline: boolean;
  onSend: (trait: string, value: TraitValue) => void;
  onBrightnessChange: (value: number) => void;
  onBrightnessCommit: (value: number) => void;
}) {
  const on = isOn(device);
  const warmth = kelvinRange(device);
  const profile = deviceProfile(device);
  const sensors = readings(device);

  const extras = writableTraits(device).filter(
    (capability) =>
      capability.trait !== TRAIT_POWER &&
      capability.trait !== TRAIT_BRIGHTNESS &&
      capability.trait !== TRAIT_COLOR_TEMP,
  );

  return (
    <div className="grid gap-5">
      <div className="flex items-center justify-between gap-4">
        <Caption>Power</Caption>
        <Switch
          checked={on}
          disabled={offline}
          onCheckedChange={(next) => onSend(TRAIT_POWER, next)}
          aria-label={`${device.name} power`}
        />
      </div>

      {dimmable ? (
        <div className="grid gap-2.5">
          <div className="flex items-baseline justify-between">
            <Caption>Brightness</Caption>
            <output className="font-mono text-small tabular">{brightness}%</output>
          </div>
          <Slider
            value={[brightness]}
            min={dimming.min}
            max={dimming.max}
            step={dimming.step}
            lit
            disabled={offline}
            className={cn(!on && "opacity-40")}
            aria-label="Brightness"
            onValueChange={([value]) => onBrightnessChange(value ?? dimming.min)}
            onValueCommit={([value]) => onBrightnessCommit(value ?? dimming.min)}
            style={{ "--lit": deviceColor(device) } as React.CSSProperties}
          />
        </div>
      ) : null}

      {warmth ? (
        <div className="grid gap-2.5">
          <Caption>Colour</Caption>
          <ColorControl
            color={deviceColor(device)}
            kelvin={numberOf(device, TRAIT_COLOR_TEMP) ?? null}
            kelvinMin={warmth.min}
            kelvinMax={warmth.max}
            showColor={hasColor(device)}
            disabled={offline}
            onWhitePreview={(kelvin) => onSend(TRAIT_COLOR_TEMP, kelvin)}
            onWhite={(kelvin) => onSend(TRAIT_COLOR_TEMP, kelvin)}
            onColor={(hex) => onSend(TRAIT_COLOR, hex)}
          />
        </div>
      ) : null}

      {extras.map((capability) => (
        <TraitControl
          key={capability.trait}
          capability={capability}
          value={device.state[capability.trait]}
          disabled={offline}
          onChange={(value) => onSend(capability.trait, value)}
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

      <div className="grid gap-5 border-t border-border pt-5">
        <RenameDevice device={device} />

        <RoomPicker device={device} />

        {device.kind === "plug" ? <AppliancePicker device={device} /> : null}

        <div className="grid gap-1 text-caption text-subtle">
          <span>{profile.form}</span>
          <span className="font-mono">{device.external_id}</span>
          {device.vendor ? <span>{device.vendor}</span> : null}
        </div>
      </div>
    </div>
  );
}

function Caption({ children }: { children: React.ReactNode }) {
  return (
    <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
      {children}
    </span>
  );
}
