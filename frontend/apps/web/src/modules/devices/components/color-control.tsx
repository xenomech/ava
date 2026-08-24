import { Slider, Tabs, TabsContent, TabsList, TabsTrigger } from "@ava/ui";
import { useState } from "react";

import { cssToHue, hueToHex } from "@/shared/lib/color";
import { useLiveSlider } from "../use-live-slider";

const KELVIN_STEP = 50;

/** Where the white slider starts on a bulb that has only ever shown colour. */
const DEFAULT_WHITE = 2700;

const SWATCHES = [
  { css: "#ff5a3c", name: "Ember" },
  { css: "#ffb347", name: "Amber" },
  { css: "#ffe28a", name: "Candle" },
  { css: "#b8ff6b", name: "Lime" },
  { css: "#4bffb5", name: "Mint" },
  { css: "#4bd6ff", name: "Lagoon" },
  { css: "#6b8cff", name: "Dusk" },
  { css: "#b46bff", name: "Violet" },
  { css: "#ff6bd6", name: "Blossom" },
  { css: "#ffffff", name: "Cool white" },
] as const;

const KELVIN_RAMP =
  "linear-gradient(to right, #ff9233, #ffbe7a, #ffe4c4, #ffffff, #d6e6ff, #a8c8ff)";
const HUE_RAMP =
  "linear-gradient(to right,#ff2d2d,#ffd400 17%,#4bff4b 33%,#00e5ff 50%,#2d6bff 67%,#b44bff 83%,#ff2d8a)";

const clamp = (value: number, low: number, high: number) => Math.min(Math.max(value, low), high);

/**
 * A bulb is either showing a colour or a white temperature, never both — the
 * hub nulls whichever one you did not set. So the tabs are the light's actual
 * mode, read from state rather than remembered from mount, and choosing one
 * applies it. A tab that only changed which slider was on screen left the bulb
 * violet while the panel claimed 2700K.
 *
 * The last white temperature is kept while the bulb is off in colour, because
 * `color_temp` comes back null and there is otherwise nothing to return to —
 * the old code fell back to the bottom of the range, so every trip through
 * colour dragged the light to its warmest setting.
 */
export function ColorControl({
  color,
  kelvin,
  kelvinMin,
  kelvinMax,
  showColor = false,
  disabled = false,
  onWhitePreview,
  onWhite,
  onColor,
}: {
  color: string;
  kelvin: number | null;
  kelvinMin: number;
  kelvinMax: number;
  showColor?: boolean;
  disabled?: boolean;
  onWhitePreview: (kelvin: number) => void;
  onWhite: (kelvin: number) => void;
  onColor: (color: string) => void;
}) {
  const [lastWhite, setLastWhite] = useState(() =>
    clamp(kelvin ?? DEFAULT_WHITE, kelvinMin, kelvinMax),
  );

  /* Adjusting state during render rather than in an effect: the slider must
     never paint a stale temperature for a frame. */
  if (kelvin !== null && kelvin !== lastWhite) setLastWhite(kelvin);

  const white = useLiveSlider(lastWhite, onWhitePreview, onWhite);
  const hue = useLiveSlider(
    Math.round((cssToHue(color) / 360) * 100),
    (percent) => onColor(hueToHex((percent / 100) * 360)),
    (percent) => onColor(hueToHex((percent / 100) * 360)),
  );

  const mode = showColor && kelvin === null ? "color" : "white";

  const choose = (next: string) => {
    if (next === mode) return;

    if (next === "white") onWhite(lastWhite);
    else onColor(hueToHex((hue.value / 100) * 360));
  };

  return (
    /* Manual activation: with Radix's default, focusing a tab selects it, so a
       click fired the command twice and arrowing between them would drive the
       bulb on the way past. */
    <Tabs
      value={mode}
      onValueChange={choose}
      activationMode="manual"
      className="grid gap-2.5"
    >
      {showColor ? (
        <TabsList className="w-full">
          <TabsTrigger value="white">White</TabsTrigger>
          <TabsTrigger value="color">Colour</TabsTrigger>
        </TabsList>
      ) : null}

      <TabsContent value="white" className="grid gap-2">
        <Slider
          value={[white.value]}
          min={kelvinMin}
          max={kelvinMax}
          step={KELVIN_STEP}
          variant="marker"
          disabled={disabled}
          aria-label="White temperature"
          aria-valuetext={`${white.value} kelvin`}
          onValueChange={([v]) => white.change(v ?? kelvinMin)}
          onValueCommit={([v]) => white.release(v ?? kelvinMin)}
          style={{ background: KELVIN_RAMP }}
        />
        <p className="font-mono text-caption text-subtle tabular">{white.value}K</p>
      </TabsContent>

      <TabsContent value="color" className={showColor ? "grid gap-2.5" : "hidden"}>
        <Slider
          value={[hue.value]}
          max={100}
          step={1}
          variant="marker"
          disabled={disabled}
          aria-label="Hue"
          /* Without onValueChange a controlled Radix thumb does not move at
             all: the slider looked broken because it was frozen at whatever
             the device last reported. */
          onValueChange={([v]) => hue.change(v ?? 0)}
          onValueCommit={([v]) => hue.release(v ?? 0)}
          style={{ background: HUE_RAMP }}
        />

        <div className="grid grid-cols-5 gap-1.5">
          {SWATCHES.map((s) => (
            <button
              key={s.css}
              type="button"
              disabled={disabled}
              aria-label={s.name}
              aria-pressed={color === s.css}
              onClick={() => onColor(s.css)}
              style={{ background: s.css }}
              className="aspect-square rounded-sm border border-border transition-transform duration-150 ease-out active:scale-95 disabled:opacity-40 aria-pressed:ring-2 aria-pressed:ring-fg aria-pressed:ring-offset-2 aria-pressed:ring-offset-surface"
            />
          ))}
        </div>
      </TabsContent>
    </Tabs>
  );
}
