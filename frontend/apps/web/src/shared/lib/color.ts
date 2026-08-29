export type Rgb = [number, number, number];
export type Hsl = { h: number; s: number; l: number };

const clamp = (value: number, low: number, high: number) => Math.min(Math.max(value, low), high);

export function rgbToCss([r, g, b]: Rgb): string {
  return `rgb(${r} ${g} ${b})`;
}

/** Pulls a colour part of the way towards another. `amount` is 0 to 1. */
export function mix(from: Rgb, towards: Rgb, amount: number): Rgb {
  return from.map((value, channel) =>
    Math.round(value + ((towards[channel] ?? value) - value) * amount),
  ) as Rgb;
}

/* parseColor sits on the drag path, so its patterns live here rather than
   being re-created per call. */
const HEX_PATTERN = /^#([0-9a-f]{3}|[0-9a-f]{6})$/i;
const HSL_PATTERN = /^hsla?\(\s*([\d.]+)(?:deg)?[\s,]+([\d.]+)%[\s,]+([\d.]+)%/i;
const RGB_PATTERN = /^rgba?\(([^)]+)\)$/i;
const RGB_SEPARATOR = /[\s,/]+/;

/**
 * Reads the two forms a device colour arrives in: a hex from a bulb that was
 * given an explicit colour, and the `rgb(r g b)` this app writes for anything
 * described by colour temperature.
 */
export function parseColor(value: string): Rgb | null {
  const text = value.trim();

  const hex = HEX_PATTERN.exec(text);

  if (hex?.[1]) {
    const digits =
      hex[1].length === 3
        ? hex[1]
            .split("")
            .map((digit) => digit + digit)
            .join("")
        : hex[1];

    return [
      parseInt(digits.slice(0, 2), 16),
      parseInt(digits.slice(2, 4), 16),
      parseInt(digits.slice(4, 6), 16),
    ];
  }

  const hsl = HSL_PATTERN.exec(text);

  if (hsl) {
    return hslToRgb({ h: Number(hsl[1]), s: Number(hsl[2]) / 100, l: Number(hsl[3]) / 100 });
  }

  const parts = RGB_PATTERN.exec(text)?.[1]?.split(RGB_SEPARATOR).filter(Boolean);

  if (!parts || parts.length < 3) return null;

  const channels = parts.slice(0, 3).map(Number);

  if (channels.some(Number.isNaN)) return null;

  return channels.map((channel) => clamp(Math.round(channel), 0, 255)) as Rgb;
}

export function rgbToHsl([r, g, b]: Rgb): Hsl {
  const red = r / 255;
  const green = g / 255;
  const blue = b / 255;

  const max = Math.max(red, green, blue);
  const min = Math.min(red, green, blue);
  const span = max - min;
  const l = (max + min) / 2;

  if (span === 0) return { h: 0, s: 0, l };

  const s = l > 0.5 ? span / (2 - max - min) : span / (max + min);

  const h =
    max === red
      ? ((green - blue) / span + (green < blue ? 6 : 0)) * 60
      : max === green
        ? ((blue - red) / span + 2) * 60
        : ((red - green) / span + 4) * 60;

  return { h, s, l };
}

export function hslToRgb({ h, s, l }: Hsl): Rgb {
  const hue = ((h % 360) + 360) % 360;
  const chroma = (1 - Math.abs(2 * l - 1)) * s;
  const second = chroma * (1 - Math.abs(((hue / 60) % 2) - 1));
  const base = l - chroma / 2;

  const [r, g, b] =
    hue < 60
      ? [chroma, second, 0]
      : hue < 120
        ? [second, chroma, 0]
        : hue < 180
          ? [0, chroma, second]
          : hue < 240
            ? [0, second, chroma]
            : hue < 300
              ? [second, 0, chroma]
              : [chroma, 0, second];

  return [
    Math.round((r + base) * 255),
    Math.round((g + base) * 255),
    Math.round((b + base) * 255),
  ] as Rgb;
}

/**
 * How warm a colour reads, as a single number. Used only for ordering, so it
 * needs to be monotonic rather than perceptually exact — and unlike hue, it
 * behaves for the near-white colours most bulbs actually sit at.
 */
export const warmth = ([r, , b]: Rgb) => r - b;
