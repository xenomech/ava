import { cn } from "@ava/ui";

import { KELVIN_MAX, KELVIN_MIN, kelvinToRgb, mix, type Rgb } from "@/shared/lib/kelvin";

/**
 * The light a room throws when its switch is flicked: it rises on the way up,
 * falls on the way down, and leaves nothing behind either way.
 *
 * Each stop starts as the room's own temperature, offset along the kelvin
 * scale, and is then pulled part of the way towards a chroma anchor. Kelvin
 * alone was the honest answer and the wrong one: between 2200K and 6100K every
 * stop is a shade of orange-to-white, so at the opacity this can afford they
 * collapsed into one indistinct smudge. The anchors put the colour back, the
 * kelvin base keeps it the room's colour, and a warm living room still sweeps
 * warmer than a cool kitchen.
 */
const STOPS: { offset: number; anchor: Rgb; amount: number; at: number }[] = [
  { offset: -700, anchor: [252, 43, 163], amount: 0.46, at: 22 },
  { offset: -200, anchor: [252, 109, 53], amount: 0.4, at: 37 },
  { offset: 500, anchor: [249, 200, 61], amount: 0.3, at: 52 },
  { offset: 1600, anchor: [194, 214, 225], amount: 0.34, at: 67 },
  { offset: 3200, anchor: [20, 78, 197], amount: 0.62, at: 82 },
];

/**
 * Fades the left and right of the plume so it reads as light rather than as a
 * lit rectangle. The top and bottom are handled in the ramp itself: an element
 * edge that cuts a gradient mid-colour draws a hard line straight across the
 * room, which is exactly the boxiness this is trying to avoid.
 */
const MASK =
  "linear-gradient(to right, rgb(0 0 0 / 0) 0%, rgb(0 0 0) 24%, rgb(0 0 0) 76%, rgb(0 0 0 / 0) 100%)";

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

  /* Going off, the cool end leads: the warm light is what drains away last. */
  const order = direction === "on" ? STOPS : [...STOPS].reverse();

  const ramp = order.flatMap((stop, index) => {
    const base = Math.min(Math.max(kelvin + stop.offset, KELVIN_MIN), KELVIN_MAX);
    const [r, g, b] = mix(kelvinToRgb(base), stop.anchor, stop.amount);
    const color = `rgb(${r} ${g} ${b})`;
    const at = STOPS[index]?.at ?? 0;

    /* Both ends run out to zero alpha so the element's own edge never shows. */
    if (index === 0) return [`rgb(${r} ${g} ${b} / 0) 0%`, `${color} ${at}%`];
    if (index === order.length - 1) return [`${color} ${at}%`, `rgb(${r} ${g} ${b} / 0) 100%`];

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
