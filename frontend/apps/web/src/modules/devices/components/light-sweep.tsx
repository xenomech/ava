import { cn } from "@ava/ui";

import { KELVIN_MAX, KELVIN_MIN, kelvinToRgb, mix, rgbToCss, type Rgb } from "@/shared/lib/kelvin";

/**
 * The band of light a room throws when its switch is flicked: it rises on the
 * way up, falls on the way down, and leaves nothing behind either way.
 *
 * Each band starts as the room's own temperature, offset a little along the
 * kelvin scale, and is then pulled part of the way towards a chroma anchor.
 * Kelvin alone was the honest answer and the wrong one: between 2200K and
 * 6100K every stop is some shade of orange-to-white, so at the opacity this
 * can afford they collapsed into one indistinct smudge. The anchors put the
 * colour back, the kelvin base keeps it the room's colour, and a warm living
 * room still sweeps warmer than a cool kitchen.
 */
const BANDS: { offset: number; anchor: Rgb; amount: number }[] = [
  { offset: -700, anchor: [252, 43, 163], amount: 0.38 },
  { offset: -200, anchor: [252, 109, 53], amount: 0.34 },
  { offset: 500, anchor: [249, 200, 61], amount: 0.22 },
  { offset: 1600, anchor: [194, 214, 225], amount: 0.28 },
  { offset: 3200, anchor: [20, 78, 197], amount: 0.55 },
];

/** Three columns, the middle one lifted, so the leading edge is uneven. */
const COLUMNS = [0, -72, 0];

export function LightSweep({
  kelvin,
  direction,
  /** Bumped on every flick. Remounts the element so the animation replays. */
  play,
}: {
  kelvin: number;
  direction: "on" | "off";
  play: number;
}) {
  if (play === 0) return null;

  const bands = BANDS.map(({ offset, anchor, amount }) => {
    const base = Math.min(Math.max(kelvin + offset, KELVIN_MIN), KELVIN_MAX);

    return rgbToCss(mix(kelvinToRgb(base), anchor, amount));
  });

  return (
    <div
      key={play}
      aria-hidden
      className={cn(
        "pointer-events-none absolute inset-y-0 left-1/2 z-base flex w-full max-w-[860px]",
        direction === "on" ? "animate-light-up" : "animate-light-down",
      )}
    >
      {COLUMNS.map((lift, column) => (
        <div
          key={column}
          className="flex flex-1 flex-col items-stretch"
          style={{ transform: `translateY(${lift}px)` }}
        >
          {bands.map((color, band) => (
            <span
              key={band}
              className="block flex-1 blur-xl"
              style={{ background: color, marginTop: band === 0 ? undefined : -12 }}
            />
          ))}
        </div>
      ))}
    </div>
  );
}
