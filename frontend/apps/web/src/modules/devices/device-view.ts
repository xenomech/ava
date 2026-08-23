import type { DeviceKind } from "@ava/ui";
import {
  TRAIT_BRIGHTNESS,
  TRAIT_COLOR,
  TRAIT_COLOR_TEMP,
  deviceProfile,
  isOn,
  numberOf,
  supports,
  traitValue,
  type DeviceDto,
} from "@ava/contracts";

import { kelvinToCss } from "@/shared/lib/kelvin";

const RESTING_KELVIN = 2700;

// How a device is drawn: which shape stands in for it, how bright it reads, and
// what colour it casts. Plain functions in a plain module — keeping them beside
// a component broke Fast Refresh for every file that imported them.

export function deviceKind(device: DeviceDto): DeviceKind {
  return deviceProfile(device).kind;
}

export function deviceLevel(device: DeviceDto): number {
  if (!isOn(device)) return 0;
  if (!supports(device, TRAIT_BRIGHTNESS)) return 100;

  return numberOf(device, TRAIT_BRIGHTNESS) ?? 100;
}

export function deviceColor(device: DeviceDto): string {
  const picked = traitValue(device, TRAIT_COLOR);
  if (typeof picked === "string" && picked !== "") return picked;

  return kelvinToCss(numberOf(device, TRAIT_COLOR_TEMP) ?? RESTING_KELVIN);
}
