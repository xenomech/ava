import { cn } from "@ava/ui";
import { useMemo } from "react";

import { hslToRgb, mix, parseColor, rgbToCss, rgbToHsl, type Rgb } from "@/shared/lib/color";

// The shape of the sweep, not its colour: every stop is a transform of the room's own light.
const SHAPE = [
  { hue: -14, shade: 0.46, at: 22 },
  { hue: -5, shade: 0.18, at: 37 },
  { hue: 0, shade: -0.1, at: 52 },
  { hue: 10, shade: 0.24, at: 67 },
  { hue: 24, shade: 0.58, at: 82 },
];

const BLACK: Rgb = [0, 0, 0];
const WHITE: Rgb = [255, 255, 255];

// Fades the sides of the plume so it reads as light rather than as a lit rectangle.
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

/** The whole gradient for one palette and direction, or null for no palette. */
function buildRamp(colors: string[], direction: "on" | "off"): string | null {
  const palette = colors.map(parseColor).filter((color): color is Rgb => color !== null);

  if (palette.length === 0) return null;

  // Going off, the cool end leads: the warm light is what drains away last.
  const shape = direction === "on" ? SHAPE : [...SHAPE].reverse();

  const ramp = shape.flatMap((step, index) => {
    const source = rgbToHsl(sample(palette, index / (SHAPE.length - 1)));
    const rotated = hslToRgb({ ...source, h: source.h + step.hue });

    const [r, g, b] =
      step.shade >= 0 ? mix(rotated, BLACK, step.shade) : mix(rotated, WHITE, -step.shade);

    const color = rgbToCss([r, g, b]);
    const at = SHAPE[index]?.at ?? 0;

    // Both ends run out to zero alpha so the element's own edge never shows.
    if (index === 0) return [`rgb(${r} ${g} ${b} / 0) 0%`, `${color} ${at}%`];
    if (index === shape.length - 1) return [`${color} ${at}%`, `rgb(${r} ${g} ${b} / 0) 100%`];

    return [`${color} ${at}%`];
  });

  return `linear-gradient(to bottom, ${ramp.join(", ")})`;
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
  // The page re-renders at pointer rate, but the ramp only changes with colours or direction.
  const ramp = useMemo(() => buildRamp(colors, direction), [colors, direction]);

  if (play === 0 || ramp === null) return null;

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
          background: ramp,
          maskImage: MASK,
          WebkitMaskImage: MASK,
        }}
      />
    </div>
  );
}
