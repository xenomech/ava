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

export function deviceKind(device: DeviceDto): DeviceKind {
  return deviceProfile(device).kind;
}

export function deviceLevel(device: DeviceDto): number {
  if (!isOn(device)) return 0;
  if (!supports(device, TRAIT_BRIGHTNESS)) return 100;

  return numberOf(device, TRAIT_BRIGHTNESS) ?? 100;
}

/**
 * What a device is doing, in the fewest words that stay true.
 *
 * A plug or a fan has no brightness, so `deviceLevel` reports 100 for them and
 * a bare percentage would claim a dimmer that is not there. Those read "On".
 */
export function deviceLabel(device: DeviceDto): string {
  if (device.status === "offline") return "Offline";
  if (!isOn(device)) return "Off";
  if (!supports(device, TRAIT_BRIGHTNESS)) return "On";

  return `${deviceLevel(device)}%`;
}

export function deviceColor(device: DeviceDto): string {
  const picked = traitValue(device, TRAIT_COLOR);
  if (typeof picked === "string" && picked !== "") return picked;

  return kelvinToCss(numberOf(device, TRAIT_COLOR_TEMP) ?? 2700);
}
