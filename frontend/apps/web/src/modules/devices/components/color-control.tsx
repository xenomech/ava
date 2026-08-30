import { MarkerSlider, cn } from "@ava/ui";

import { parseColor } from "@/shared/lib/color";
import { kelvinToCss } from "@/shared/lib/kelvin";
import { ColourPad, tintCss, tintHex, tintOf, type Tint } from "./colour-pad";
import { useLiveSlider } from "../use-live-slider";

const KELVIN_STEP = 50;

/** Where the whites start on a bulb that has only ever shown colour. */
const DEFAULT_WHITE = 2700;

/**
 * The looks a room is actually set to, as positions on the pad.
 *
 * Written in the pad's own coordinates rather than as hex, so every swatch is
 * somewhere the thumb can land. Hex carries a lightness the pad does not model,
 * and a palette whose own colours its control cannot reproduce is a palette
 * that argues with you: tapping one and then touching the pad used to hand the
 * bulb a different colour than the one you had picked.
 */
const WHITES = [2200, 2700, 3500, 4600, 6000].map((kelvin) => ({
  kelvin,
  css: kelvinToCss(kelvin),
}));

/* The css and its parsed channels are fixed, so they are computed once here
   rather than per swatch per render while the pad is being dragged. */
const COLOURS = (
  [
    { hue: 8, whiteness: 0.08, name: "Ember" },
    { hue: 32, whiteness: 0.05, name: "Amber" },
    { hue: 44, whiteness: 0.3, name: "Candle" },
    { hue: 28, whiteness: 0.62, name: "Sand" },
    { hue: 92, whiteness: 0.22, name: "Lime" },
    { hue: 156, whiteness: 0.3, name: "Mint" },
    { hue: 194, whiteness: 0.15, name: "Lagoon" },
    { hue: 224, whiteness: 0.28, name: "Dusk" },
    { hue: 274, whiteness: 0.32, name: "Violet" },
    { hue: 320, whiteness: 0.28, name: "Blossom" },
  ] as const
).map(({ hue, whiteness, name }) => {
  const tint: Tint = { mode: "colour", hue, whiteness };
  const css = tintCss(tint);
  return { name, tint, css, rgb: parseColor(css) };
});

const KELVIN_RAMP =
  "linear-gradient(to right, #ff9233, #ffbe7a, #ffe4c4, #ffffff, #d6e6ff, #a8c8ff)";

const clamp = (value: number, low: number, high: number) => Math.min(Math.max(value, low), high);

/**
 * The warmth track for a bulb that only does tunable white — a pad whose colour
 * half is unreachable would be a lie about the hardware.
 */
export function WhiteControl({
  kelvin,
  kelvinMin,
  kelvinMax,
  disabled = false,
  onWhite,
}: {
  kelvin: number | null;
  kelvinMin: number;
  kelvinMax: number;
  disabled?: boolean;
  onWhite: (kelvin: number) => void;
}) {
  const white = useLiveSlider(
    clamp(kelvin ?? DEFAULT_WHITE, kelvinMin, kelvinMax),
    onWhite,
    onWhite,
  );

  return (
    <div className="grid gap-2">
      <MarkerSlider
        value={[white.value]}
        min={kelvinMin}
        max={kelvinMax}
        step={KELVIN_STEP}
        disabled={disabled}
        aria-label="White temperature"
        aria-valuetext={`${white.value} kelvin`}
        onValueChange={([value]) => white.change(value ?? kelvinMin)}
        onValueCommit={([value]) => white.release(value ?? kelvinMin)}
        style={{ background: KELVIN_RAMP }}
      />
      <p className="font-mono text-caption text-subtle tabular">{white.value}K</p>
    </div>
  );
}

/**
 * Everything a light can be, on one surface.
 *
 * There is no White tab and no Colour tab. A bulb shows a colour or a
 * temperature and never both, and the hub clears whichever one you did not set
 * — so the mode was never really a setting, it was a fact about where the light
 * already was. Every colour bug in this app has come from the seam between two
 * tabs claiming otherwise: a tab that only swapped which slider was on screen
 * left the bulb violet while the panel read 2700K, and a swatch that asked for
 * RGB white left the panel reading Colour while the lamp looked white.
 *
 * The pad has no opinion to be wrong about. The app reads which trait to send
 * from where the thumb is.
 */
export function ColorControl({
  color,
  kelvin,
  kelvinMin,
  kelvinMax,
  disabled = false,
  onWhite,
  onColor,
}: {
  color: string;
  kelvin: number | null;
  kelvinMin: number;
  kelvinMax: number;
  disabled?: boolean;
  onWhite: (kelvin: number) => void;
  onColor: (color: string) => void;
}) {
  const send = (tint: Tint) => {
    if (tint.mode === "white") onWhite(Math.round(tint.kelvin));
    else onColor(tintHex(tint));
  };

  const settled = tintOf(color, kelvin, kelvinMin, kelvinMax);
  const pad = useLiveSlider<Tint>(settled, send, send);

  /* Parsed once per render, compared against each swatch's precomputed rgb. */
  const currentRgb = kelvin === null ? parseColor(color) : null;

  return (
    <div className="grid gap-2.5">
      <ColourPad
        tint={pad.value}
        kelvinMin={kelvinMin}
        kelvinMax={kelvinMax}
        disabled={disabled}
        onPreview={pad.change}
        onCommit={pad.release}
        onCancel={pad.reset}
      />

      <p className="flex items-baseline justify-between gap-4 font-mono text-caption text-subtle">
        <span>{pad.value.mode === "white" ? `${Math.round(pad.value.kelvin)}K` : "Colour"}</span>
        <span className="tabular">
          {pad.value.mode === "white" ? "white" : tintHex(pad.value).toUpperCase()}
        </span>
      </p>

      <div className="grid grid-cols-5 gap-1.5">
        {WHITES.map(({ kelvin: degrees, css }) => (
          <Swatch
            key={degrees}
            label={`${degrees}K white`}
            css={css}
            pressed={pad.value.mode === "white" && Math.abs(pad.value.kelvin - degrees) < 60}
            disabled={disabled}
            onPick={() => pad.release({ mode: "white", kelvin: degrees })}
          />
        ))}

        {COLOURS.map((swatch) => (
          <Swatch
            key={swatch.name}
            label={swatch.name}
            css={swatch.css}
            pressed={sameRgb(currentRgb, swatch.rgb)}
            disabled={disabled}
            onPick={() => pad.release(swatch.tint)}
          />
        ))}
      </div>
    </div>
  );
}

function Swatch({
  label,
  css,
  pressed,
  disabled,
  onPick,
}: {
  label: string;
  css: string;
  pressed: boolean;
  disabled: boolean;
  onPick: () => void;
}) {
  return (
    <button
      type="button"
      disabled={disabled}
      aria-label={label}
      aria-pressed={pressed}
      onClick={onPick}
      style={{ background: css }}
      className={cn(
        "aspect-square rounded-sm border border-border transition-transform duration-150",
        "ease-out active:scale-95 disabled:opacity-40",
        "aria-pressed:ring-2 aria-pressed:ring-fg aria-pressed:ring-offset-2",
        "aria-pressed:ring-offset-surface",
      )}
    />
  );
}

/**
 * Whether two colours are the same lamp setting.
 *
 * Compared as numbers, not as strings: the same colour arrives as `#FFB347`
 * from one place and `rgb(255 179 71)` from another, so a string comparison
 * meant a swatch you had just tapped never looked tapped.
 */
function sameRgb(
  a: readonly number[] | null | undefined,
  b: readonly number[] | null | undefined,
): boolean {
  if (!a || !b) return false;

  /* A hair of tolerance, because a bulb rounds what it is given and reports the
     rounded value back. */
  return a.every((channel, at) => Math.abs(channel - (b[at] ?? 0)) <= 2);
}
