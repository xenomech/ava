import { Slider, Tabs, TabsContent, TabsList, TabsTrigger } from "@ava/ui";
import { useState } from "react";

import { cssToTint, parseColor, tintToHex, type Tint } from "@/shared/lib/color";
import { useLiveSlider } from "../use-live-slider";

const KELVIN_STEP = 50;

/** Where the white slider starts on a bulb that has only ever shown colour. */
const DEFAULT_WHITE = 2700;

/**
 * Ten colours a room actually gets set to, warm to cool.
 *
 * Written as the two numbers the sliders work in rather than as hex, so every
 * swatch is somewhere the sliders can actually land. They were hex before, and
 * hex carries a lightness the picker does not model — so tapping Lagoon and
 * then so much as changing tabs handed the bulb `#00c5ff` instead of the
 * `#4bd6ff` you had chosen. A palette whose own colours it cannot reproduce is
 * a palette that argues with you.
 *
 * "Cool white" used to close the row and was a different trap: it asks for RGB
 * white, which is not the tunable white on the other tab, and left the panel
 * showing Colour while the lamp looked like it was showing White. Saturation
 * reaches it now for anyone who wants it.
 */
const SWATCHES = [
  { hue: 12, saturation: 90, name: "Ember" },
  { hue: 32, saturation: 95, name: "Amber" },
  { hue: 44, saturation: 72, name: "Candle" },
  { hue: 28, saturation: 34, name: "Sand" },
  { hue: 92, saturation: 78, name: "Lime" },
  { hue: 156, saturation: 70, name: "Mint" },
  { hue: 194, saturation: 85, name: "Lagoon" },
  { hue: 224, saturation: 72, name: "Dusk" },
  { hue: 274, saturation: 68, name: "Violet" },
  { hue: 320, saturation: 72, name: "Blossom" },
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

  const mode = showColor && kelvin === null ? "color" : "white";

  /* The colour the lamp is holding, kept while it shows white — otherwise the
     sliders read the kelvin ramp's own orange, and coming back to Colour landed
     on a shade nobody had chosen. */
  const [lastTint, setLastTint] = useState(() => cssToTint(color));
  const shown = mode === "color" ? cssToTint(color) : lastTint;

  if (
    mode === "color" &&
    (shown.hue !== lastTint.hue || shown.saturation !== lastTint.saturation)
  ) {
    setLastTint(shown);
  }

  const send = (next: Partial<Tint>) => {
    const tint = { ...shown, ...next };

    setLastTint(tint);
    onColor(tintToHex(tint));
  };

  const hue = useLiveSlider(
    Math.round(shown.hue),
    (value) => send({ hue: value }),
    (value) => send({ hue: value }),
  );

  const saturation = useLiveSlider(
    shown.saturation,
    (value) => send({ saturation: value }),
    (value) => send({ saturation: value }),
  );

  const choose = (next: string) => {
    if (next === mode) return;

    if (next === "white") onWhite(lastWhite);
    else onColor(tintToHex(lastTint));
  };

  return (
    /* Manual activation: with Radix's default, focusing a tab selects it, so a
       click fired the command twice and arrowing between them would drive the
       bulb on the way past. */
    <Tabs value={mode} onValueChange={choose} activationMode="manual" className="grid gap-2.5">
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
          max={360}
          step={1}
          variant="marker"
          disabled={disabled}
          aria-label="Hue"
          aria-valuetext={`${hue.value} degrees`}
          /* Without onValueChange a controlled Radix thumb does not move at
             all: the slider looked broken because it was frozen at whatever
             the device last reported. */
          onValueChange={([v]) => hue.change(v ?? 0)}
          onValueCommit={([v]) => hue.release(v ?? 0)}
          style={{ background: HUE_RAMP }}
        />

        {/* The other half of the wheel. Hue says which colour, this says how
            much of it, and between them every shade the swatches offer becomes
            reachable by hand. The track is drawn in the hue you are on, so it
            previews its own outcome. */}
        <Slider
          value={[saturation.value]}
          max={100}
          step={1}
          variant="marker"
          disabled={disabled}
          aria-label="Saturation"
          aria-valuetext={`${saturation.value} percent`}
          onValueChange={([v]) => saturation.change(v ?? 0)}
          onValueCommit={([v]) => saturation.release(v ?? 0)}
          style={{
            background: `linear-gradient(to right, ${tintToHex({ hue: hue.value, saturation: 0 })}, ${tintToHex({ hue: hue.value, saturation: 100 })})`,
          }}
        />

        <div className="grid grid-cols-5 gap-1.5">
          {SWATCHES.map((s) => (
            <button
              key={s.name}
              type="button"
              disabled={disabled}
              aria-label={s.name}
              aria-pressed={sameColor(color, tintToHex(s))}
              onClick={() => send(s)}
              style={{ background: tintToHex(s) }}
              className="aspect-square rounded-sm border border-border transition-transform duration-150 ease-out active:scale-95 disabled:opacity-40 aria-pressed:ring-2 aria-pressed:ring-fg aria-pressed:ring-offset-2 aria-pressed:ring-offset-surface"
            />
          ))}
        </div>
      </TabsContent>
    </Tabs>
  );
}

/**
 * Whether two colours are the same lamp setting.
 *
 * Compared as numbers, not as strings: the same colour arrives as `#FFB347`
 * from one place and `rgb(255 179 71)` from another, so a string comparison
 * meant a swatch you had just tapped never looked tapped.
 */
function sameColor(left: string, right: string): boolean {
  const a = parseColor(left);
  const b = parseColor(right);

  if (!a || !b) return false;

  return a[0] === b[0] && a[1] === b[1] && a[2] === b[2];
}
