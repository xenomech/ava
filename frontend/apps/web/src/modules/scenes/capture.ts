import {
  TRAIT_BRIGHTNESS,
  TRAIT_COLOR,
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  isOn,
  numberOf,
  supports,
  traitValue,
  type DeviceDto,
  type SceneDto,
  type SceneTargetDto,
} from "@ava/contracts";

/**
 * The room, frozen: what to save when someone says "remember this".
 *
 * Only what a person would recognise on the switch — whether each thing is on,
 * and if it is, how bright and what colour. The rest of a device's traits are
 * either readings, which cannot be set, or settings nobody thinks of as part of
 * "how the room looks".
 *
 * A device that is off is stored as off and nothing more. Its brightness while
 * dark is not a fact about the room, and replaying it would make the light
 * flare at the old level for the moment before the power target lands.
 */
export function capture(devices: DeviceDto[]): SceneTargetDto[] {
  return devices.flatMap((device) => {
    if (device.room_id === undefined || !supports(device, TRAIT_POWER)) return [];

    const on = isOn(device);
    const targets: SceneTargetDto[] = [{ device_id: device.id, trait: TRAIT_POWER, value: on }];

    if (!on) return targets;

    const brightness = numberOf(device, TRAIT_BRIGHTNESS);
    if (supports(device, TRAIT_BRIGHTNESS) && brightness !== undefined) {
      targets.push({ device_id: device.id, trait: TRAIT_BRIGHTNESS, value: brightness });
    }

    /* A light holds a colour or a temperature, never both — the hub clears
       whichever one you did not set — so the scene stores whichever it is
       actually showing. Storing both would mean replaying one and immediately
       contradicting it with the other. */
    const color = traitValue(device, TRAIT_COLOR);
    const kelvin = numberOf(device, TRAIT_COLOR_TEMP);

    if (typeof color === "string" && supports(device, TRAIT_COLOR)) {
      targets.push({ device_id: device.id, trait: TRAIT_COLOR, value: color });
    } else if (kelvin !== undefined && supports(device, TRAIT_COLOR_TEMP)) {
      targets.push({ device_id: device.id, trait: TRAIT_COLOR_TEMP, value: kelvin });
    }

    return targets;
  });
}

/**
 * Whether the room is, right now, doing what this scene says.
 *
 * Numbers are compared loosely. A bulb asked for 60% reports back 60% most of
 * the time and 61% sometimes, and a scene that can never light up because of
 * rounding is worse than no indicator at all. Two percent of the value covers
 * that drift without letting genuinely different settings pass — at 2700K it
 * allows 54 degrees, well inside one step of the warmth slider.
 */
export function matches(scene: SceneDto, devices: DeviceDto[]): boolean {
  const known = new Map(devices.map((device) => [device.id, device]));

  const live = scene.targets.filter((target) => known.has(target.device_id));
  if (live.length === 0) return false;

  return live.every((target) => {
    const current = traitValue(known.get(target.device_id), target.trait);

    if (typeof target.value === "number") {
      return (
        typeof current === "number" &&
        Math.abs(current - target.value) <= Math.max(1, Math.abs(target.value) * 0.02)
      );
    }

    if (typeof target.value === "string" && typeof current === "string") {
      return current.toLowerCase() === target.value.toLowerCase();
    }

    return current === target.value;
  });
}

/** "Table lamp · 60% · 2700K", for the list of what is about to be saved. */
export function describe(device: DeviceDto, targets: SceneTargetDto[]): string {
  const mine = targets.filter((target) => target.device_id === device.id);
  const power = mine.find((target) => target.trait === TRAIT_POWER);

  if (power?.value !== true) return "off";

  const parts = ["on"];

  for (const target of mine) {
    if (target.trait === TRAIT_BRIGHTNESS && typeof target.value === "number") {
      parts.push(`${Math.round(target.value)}%`);
    }

    if (target.trait === TRAIT_COLOR_TEMP && typeof target.value === "number") {
      parts.push(`${Math.round(target.value)}K`);
    }

    if (target.trait === TRAIT_COLOR && typeof target.value === "string") {
      parts.push(target.value.toUpperCase());
    }
  }

  return parts.join(" · ");
}
