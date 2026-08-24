export const KELVIN_MIN = 1800;
export const KELVIN_MAX = 6500;

export type Rgb = [number, number, number];

export function kelvinToRgb(kelvin: number): Rgb {
  const t = Math.min(Math.max(kelvin, KELVIN_MIN), KELVIN_MAX) / 100;

  const r = t <= 66 ? 255 : 329.7 * Math.pow(t - 60, -0.1332);
  const g = t <= 66 ? 99.47 * Math.log(t) - 161.12 : 288.12 * Math.pow(t - 60, -0.0755);
  const b = t >= 66 ? 255 : t <= 19 ? 0 : 138.52 * Math.log(t - 10) - 305.04;

  const clamp = (v: number) => Math.round(Math.min(Math.max(v, 0), 255));

  return [clamp(r), clamp(g), clamp(b)];
}

export function kelvinToCss(kelvin: number): string {
  const [r, g, b] = kelvinToRgb(kelvin);

  return `rgb(${r} ${g} ${b})`;
}

/** Pulls a colour part of the way towards another. `amount` is 0 to 1. */
export function mix(from: Rgb, towards: Rgb, amount: number): Rgb {
  return from.map((value, channel) =>
    Math.round(value + ((towards[channel] ?? value) - value) * amount),
  ) as Rgb;
}

export function rgbToCss([r, g, b]: Rgb): string {
  return `rgb(${r} ${g} ${b})`;
}

export const percentToKelvin = (percent: number) =>
  Math.round(KELVIN_MIN + (percent / 100) * (KELVIN_MAX - KELVIN_MIN));

export const kelvinToPercent = (kelvin: number) =>
  ((kelvin - KELVIN_MIN) / (KELVIN_MAX - KELVIN_MIN)) * 100;
