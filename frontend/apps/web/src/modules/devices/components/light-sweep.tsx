import { cn } from "@ava/ui";

import { hslToRgb, mix, parseColor, rgbToCss, rgbToHsl, type Rgb } from "@/shared/lib/color";

/**
 * The light a room throws when its switch is flicked: it rises on the way up,
 * falls on the way down, and leaves nothing behind either way.
 *
 * This is the shape of the ramp, not its colour. Every stop is a transform of a
 * colour the room actually holds: the hue rotates a little across the sweep,
 * the middle sits at the bulb's own colour, and the two ends fall away into it.
 *
 * `shade` is a move towards black, or towards white when negative. Doing the
 * darkening this way rather than by cutting HSL lightness is what keeps two
 * rooms apart: a 5000K white sits at lightness 0.9 with saturation near 1, so
 * halving its lightness turns it into a vivid red indistinguishable from an
 * actual red bulb. Mixing towards black instead keeps a pale colour pale and a
 * saturated one saturated.
 *
 * An earlier version mixed towards fixed pink/orange/blue anchors. It looked
 * right in one room and like somebody else's palette in every other, because at
 * those mix amounts the anchor, not the bulb, was deciding the colour. Nothing
 * here names a colour, so a red lamp sweeps red and a 2200K filament sweeps
 * amber without either being special-cased.
 */
const SHAPE = [
  { hue: -14, shade: 0.46, at: 22 },
  { hue: -5, shade: 0.18, at: 37 },
  { hue: 0, shade: -0.1, at: 52 },
  { hue: 10, shade: 0.24, at: 67 },
  { hue: 24, shade: 0.58, at: 82 },
];

const BLACK: Rgb = [0, 0, 0];
const WHITE: Rgb = [255, 255, 255];

/**
 * Fades the left and right of the plume so it reads as light rather than as a
 * lit rectangle. The top and bottom are handled in the ramp itself: an element
 * edge that cuts a gradient mid-colour draws a hard line straight across the
 * room, which is exactly the boxiness this is trying to avoid.
 */
const MASK =
  "linear-gradient(to right, rgb(0 0 0 / 0) 0%, rgb(0 0 0) 24%, rgb(0 0 0) 76%, rgb(0 0 0 / 0) 100%)";

const clamp = (value: number, low: number, high: number) => Math.min(Math.max(value, low), high);

/** Reads the palette at a position from 0 to 1, blending between neighbours. */
function sample(palette: Rgb[], position: number): Rgb {
  const first = palette[0] ?? [255, 255, 255];
  if (palette.length === 1) return first;

  const span = (palette.length - 1) * clamp(position, 0, 1);
  const index = Math.min(Math.floor(span), palette.length - 2);

  return mix(palette[index] ?? first, palette[index + 1] ?? first, span - index);
}

export function LightSweep({
  colors,
  direction,
  /** Bumped on every flick. Remounts the element so the animation replays. */
  play,
}: {
  /** The room's own colours, warmest first. */
  colors: string[];
  direction: "on" | "off";
  play: number;
}) {
  if (play === 0) return null;

  const palette = colors.map(parseColor).filter((color): color is Rgb => color !== null);

  if (palette.length === 0) return null;

  /* Going off, the cool end leads: the warm light is what drains away last. */
  const shape = direction === "on" ? SHAPE : [...SHAPE].reverse();

  const ramp = shape.flatMap((step, index) => {
    const source = rgbToHsl(sample(palette, index / (SHAPE.length - 1)));
    const rotated = hslToRgb({ ...source, h: source.h + step.hue });

    const [r, g, b] =
      step.shade >= 0 ? mix(rotated, BLACK, step.shade) : mix(rotated, WHITE, -step.shade);

    const color = rgbToCss([r, g, b]);
    const at = SHAPE[index]?.at ?? 0;

    /* Both ends run out to zero alpha so the element's own edge never shows. */
    if (index === 0) return [`rgb(${r} ${g} ${b} / 0) 0%`, `${color} ${at}%`];
    if (index === shape.length - 1) return [`${color} ${at}%`, `rgb(${r} ${g} ${b} / 0) 100%`];

    return [`${color} ${at}%`];
  });

  return (
    <div
      key={play}
      aria-hidden
      className={cn(
        "pointer-events-none absolute inset-y-0 left-1/2 z-base w-full max-w-[860px]",
        direction === "on" ? "animate-light-up" : "animate-light-down",
      )}
    >
      <div
        className="size-full blur-2xl"
        style={{
          background: `linear-gradient(to bottom, ${ramp.join(", ")})`,
          maskImage: MASK,
          WebkitMaskImage: MASK,
        }}
      />
    </div>
  );
}
