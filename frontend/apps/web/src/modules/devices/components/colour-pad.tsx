import { cn } from "@ava/ui";
import { useRef, type KeyboardEvent, type PointerEvent } from "react";

import { hslToRgb, rgbToCss, rgbToHsl, parseColor } from "@/shared/lib/color";
import { kelvinToCss, kelvinToRgb } from "@/shared/lib/kelvin";

/**
 * Where the colour runs out and the whites begin, as a fraction of the pad.
 *
 * The band is not a second control. It is the bottom edge of this one: drag
 * down through the pastels and the colour thins until there is none left, which
 * is exactly where tunable white starts. A light is either showing a colour or
 * a temperature and never both, and the honest way to say that is a place you
 * can see rather than a tab you have to remember.
 */
const WHITE_EDGE = 0.78;

/**
 * How far towards white the bottom of the colour region gets.
 *
 * Matched exactly by the overlay gradient below, because the pad has to be a
 * picture of what it will do. Draining colour at a fixed lightness heads for
 * mid grey instead — a colour no lamp is ever asked for, and a dead end right
 * where the whites need to begin.
 */
const WASH = 0.92;

export type Tint =
  | { mode: "colour"; hue: number; whiteness: number }
  | { mode: "white"; kelvin: number };

const clamp = (value: number, low: number, high: number) => Math.min(Math.max(value, low), high);

/** The colour a position produces. The gradients are built from this. */
export function tintCss(tint: Tint): string {
  if (tint.mode === "white") return kelvinToCss(tint.kelvin);

  const base = hslToRgb({ h: tint.hue, s: 1, l: 0.5 });
  const towards = clamp(tint.whiteness, 0, 1) * WASH;

  return rgbToCss(
    base.map((channel) => Math.round(channel + (255 - channel) * towards)) as [
      number,
      number,
      number,
    ],
  );
}

/** The hex a bulb wants, for the colour half. */
export function tintHex(tint: Tint): string {
  const css = tintCss(tint);
  const rgb = parseColor(css) ?? [255, 255, 255];

  return `#${rgb.map((channel) => channel.toString(16).padStart(2, "0")).join("")}`;
}

/** Where a device's current colour sits on the pad. */
export function tintOf(color: string, kelvin: number | null, min: number, max: number): Tint {
  if (kelvin !== null) return { mode: "white", kelvin: clamp(kelvin, min, max) };

  const rgb = parseColor(color);
  if (!rgb) return { mode: "colour", hue: 30, whiteness: 0 };

  const { h, s } = rgbToHsl(rgb);

  /* Saturation read back as the distance already travelled towards white, so a
     colour set here lands the thumb where it was left. */
  return { mode: "colour", hue: h, whiteness: clamp((1 - s) / WASH, 0, 1) };
}

/**
 * One surface for every colour a bulb can make.
 *
 * Hue across, how much of it down, whites along the bottom. There is no mode
 * switch: the app reads which trait to send from where the thumb is, which is
 * the whole point — every colour bug in this app has come from the seam between
 * two tabs, and a seam you can see is one nobody has to hold in their head.
 */
export function ColourPad({
  tint,
  kelvinMin,
  kelvinMax,
  disabled = false,
  onPreview,
  onCommit,
}: {
  tint: Tint;
  kelvinMin: number;
  kelvinMax: number;
  disabled?: boolean;
  onPreview: (tint: Tint) => void;
  onCommit: (tint: Tint) => void;
}) {
  const pad = useRef<HTMLDivElement>(null);
  const dragging = useRef(false);

  const kelvinAt = (x: number) => Math.round(kelvinMin + x * (kelvinMax - kelvinMin));

  const at = (event: PointerEvent<HTMLDivElement>): Tint => {
    const box = pad.current?.getBoundingClientRect();
    if (!box) return tint;

    const x = clamp((event.clientX - box.left) / box.width, 0, 1);
    const y = clamp((event.clientY - box.top) / box.height, 0, 1);

    if (y > WHITE_EDGE) return { mode: "white", kelvin: kelvinAt(x) };

    return { mode: "colour", hue: x * 360, whiteness: clamp(y / WHITE_EDGE, 0, 1) };
  };

  const down = (event: PointerEvent<HTMLDivElement>) => {
    if (disabled) return;

    event.currentTarget.setPointerCapture(event.pointerId);
    dragging.current = true;
    onPreview(at(event));
  };

  const move = (event: PointerEvent<HTMLDivElement>) => {
    if (!dragging.current) return;

    onPreview(at(event));
  };

  const up = (event: PointerEvent<HTMLDivElement>) => {
    if (!dragging.current) return;

    dragging.current = false;
    onCommit(at(event));
  };

  /* Arrows move by a usable step rather than a pixel: across changes the hue or
     the temperature, down walks towards white and then into the band. */
  const key = (event: KeyboardEvent<HTMLDivElement>) => {
    if (disabled) return;

    const step = event.shiftKey ? 4 : 1;
    let next: Tint | null = null;

    if (tint.mode === "white") {
      const span = (kelvinMax - kelvinMin) / 40;

      if (event.key === "ArrowLeft")
        next = { mode: "white", kelvin: clamp(tint.kelvin - span * step, kelvinMin, kelvinMax) };
      if (event.key === "ArrowRight")
        next = { mode: "white", kelvin: clamp(tint.kelvin + span * step, kelvinMin, kelvinMax) };
      if (event.key === "ArrowUp") next = { mode: "colour", hue: 30, whiteness: 1 };
    } else {
      if (event.key === "ArrowLeft") next = { ...tint, hue: (tint.hue - 6 * step + 360) % 360 };
      if (event.key === "ArrowRight") next = { ...tint, hue: (tint.hue + 6 * step) % 360 };
      if (event.key === "ArrowUp")
        next = { ...tint, whiteness: clamp(tint.whiteness - 0.06 * step, 0, 1) };
      if (event.key === "ArrowDown") {
        next =
          tint.whiteness >= 1
            ? { mode: "white", kelvin: kelvinAt(0.5) }
            : { ...tint, whiteness: clamp(tint.whiteness + 0.06 * step, 0, 1) };
      }
    }

    if (!next) return;

    event.preventDefault();
    onCommit(next);
  };

  const position =
    tint.mode === "white"
      ? {
          left: `${((tint.kelvin - kelvinMin) / (kelvinMax - kelvinMin)) * 100}%`,
          top: `${(WHITE_EDGE + (1 - WHITE_EDGE) / 2) * 100}%`,
        }
      : { left: `${(tint.hue / 360) * 100}%`, top: `${tint.whiteness * WHITE_EDGE * 100}%` };

  return (
    /* The rule's list of interactive roles does not include `application`,
       which is the one role that genuinely means "this widget takes its own
       keys". There is no element or role for two continuous values on a
       surface, and the accessible path is the labelled swatch grid below. */
    // eslint-disable-next-line jsx-a11y/no-noninteractive-element-interactions
    <div
      ref={pad}
      /* `application`, because this widget handles its own arrow keys and there
         is no role that means "two continuous values on a surface". Assistive
         technology is not left with only this: the swatch grid underneath is a
         complete set of labelled buttons covering every look the pad is for. */
      role="application"
      aria-roledescription="Colour pad"
      tabIndex={disabled ? -1 : 0}
      aria-label="Colour"
      aria-disabled={disabled}
      onPointerDown={down}
      onPointerMove={move}
      onPointerUp={up}
      onPointerCancel={up}
      onKeyDown={key}
      className={cn(
        "relative h-[168px] w-full touch-none select-none overflow-hidden rounded-lg",
        "border border-border outline-none",
        "focus-visible:ring-2 focus-visible:ring-fg focus-visible:ring-offset-2",
        "focus-visible:ring-offset-surface",
        disabled ? "cursor-not-allowed opacity-40" : "cursor-pointer",
      )}
      style={{
        /* The wash is the same maths `tintCss` applies, expressed as a gradient
           — white at `WASH` alpha by `WHITE_EDGE` — so the thumb always sits on
           the colour it is about to send. */
        backgroundImage: [
          `linear-gradient(to bottom, rgb(255 255 255 / 0) 0%, rgb(255 255 255 / ${WASH}) ${WHITE_EDGE * 100}%)`,
          hueRamp(),
        ].join(","),
      }}
    >
      <span
        aria-hidden
        className="absolute inset-x-0 bottom-0 border-t border-white/30"
        style={{
          height: `${(1 - WHITE_EDGE) * 100}%`,
          backgroundImage: whiteRamp(kelvinMin, kelvinMax),
        }}
      />

      <span
        aria-hidden
        /* Right-hand side: the thumb spends most of its life at the warm end of
           the band, which is where this used to sit. */
        className="pointer-events-none absolute font-mono text-[9px] uppercase tracking-caps text-black/45"
        style={{ right: 9, bottom: 5 }}
      >
        Whites
      </span>

      <span
        aria-hidden
        className={cn(
          "pointer-events-none absolute -ml-[11px] -mt-[11px] size-[22px] rounded-full",
          "border-[2.5px] border-white",
          "shadow-[0_2px_10px_rgb(0_0_0/0.65)]",
        )}
        style={position}
      />

      <output className="sr-only">
        {tint.mode === "white"
          ? `${Math.round(tint.kelvin)} kelvin`
          : `hue ${Math.round(tint.hue)} degrees`}
      </output>
    </div>
  );
}

/** Sampled from the same function the thumb reads, so the two cannot disagree. */
function hueRamp(steps = 48): string {
  const stops = Array.from({ length: steps + 1 }, (_, index) => {
    const hue = (index / steps) * 360;

    return `${rgbToCss(hslToRgb({ h: hue, s: 1, l: 0.5 }))} ${((index / steps) * 100).toFixed(1)}%`;
  });

  return `linear-gradient(to right, ${stops.join(",")})`;
}

function whiteRamp(min: number, max: number, steps = 24): string {
  const stops = Array.from({ length: steps + 1 }, (_, index) => {
    const kelvin = min + (index / steps) * (max - min);

    return `${rgbToCss(kelvinToRgb(kelvin))} ${((index / steps) * 100).toFixed(1)}%`;
  });

  return `linear-gradient(to right, ${stops.join(",")})`;
}
