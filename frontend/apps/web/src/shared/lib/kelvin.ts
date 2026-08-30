import type { Rgb } from "./color";

export const KELVIN_MIN = 1800;
export const KELVIN_MAX = 6500;

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
