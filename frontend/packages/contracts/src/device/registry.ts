import {
  TRAIT_BRIGHTNESS,
  TRAIT_COLOR,
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  type CapabilityDto,
  type DeviceDto,
  type TraitValue,
} from "./device.dto";

export const DEVICE_FORMS = [
  "bulb",
  "tube",
  "strip",
  "lamp",
  "plug",
  "sensor",
  "fan",
  "heater",
  "speaker",
] as const;

export type DeviceForm = (typeof DEVICE_FORMS)[number];

export type DeviceProfile = {
  kind: DeviceForm;
  form: string;
};

const LIGHT_FORMS = new Set<DeviceForm>(["bulb", "tube", "strip", "lamp", "heater"]);

export function emitsLight(form: DeviceForm) {
  return LIGHT_FORMS.has(form);
}

const MODELS: Record<string, DeviceProfile> = {
  ESP25_SHTW_01: { kind: "tube", form: "Batten" },
};

const APPLIANCES: Record<string, DeviceProfile> = {
  lamp: { kind: "lamp", form: "Table lamp" },
  bulb: { kind: "bulb", form: "Bulb" },
  tube: { kind: "tube", form: "Batten" },
  strip: { kind: "strip", form: "Light strip" },
  fan: { kind: "fan", form: "Fan" },
  heater: { kind: "heater", form: "Heater" },
  speaker: { kind: "speaker", form: "Speaker" },
};

export const APPLIANCE_CHOICES = Object.entries(APPLIANCES).map(([value, p]) => ({
  value,
  label: p.form,
}));

const KINDS: Record<string, DeviceProfile> = {
  bulb: { kind: "bulb", form: "Bulb" },
  tube: { kind: "tube", form: "Batten" },
  strip: { kind: "strip", form: "Light strip" },
  lamp: { kind: "lamp", form: "Lamp" },
  plug: { kind: "plug", form: "Plug" },
  sensor: { kind: "sensor", form: "Sensor" },
  fan: { kind: "fan", form: "Fan" },
  heater: { kind: "heater", form: "Heater" },
  speaker: { kind: "speaker", form: "Speaker" },
};

const UNKNOWN: DeviceProfile = { kind: "bulb", form: "Light" };

export function deviceProfile(device: DeviceDto): DeviceProfile {
  const chosen = device.appliance && APPLIANCES[device.appliance];
  if (chosen) return chosen;

  const model = device.model;
  if (model && MODELS[model]) return MODELS[model];

  return KINDS[device.kind] ?? UNKNOWN;
}

export type Range = { min: number; max: number; step: number; unit: string };

export function writableTraits(device: DeviceDto | undefined): CapabilityDto[] {
  return (device?.capabilities ?? []).filter((capability) => capability.access === "rw");
}

export function readings(device: DeviceDto | undefined): CapabilityDto[] {
  return (device?.capabilities ?? []).filter((capability) => capability.access === "r");
}

export function capabilityFor(
  device: DeviceDto | undefined,
  trait: string,
): CapabilityDto | undefined {
  return device?.capabilities.find((capability) => capability.trait === trait);
}

export function supports(device: DeviceDto | undefined, trait: string): boolean {
  return capabilityFor(device, trait) !== undefined;
}

export function rangeOf(capability: CapabilityDto | undefined): Range | null {
  if (!capability || capability.kind !== "number") return null;

  return {
    min: capability.min ?? 0,
    max: capability.max ?? 100,
    step: capability.step ?? 1,
    unit: capability.unit ?? "",
  };
}

export function traitValue(device: DeviceDto | undefined, trait: string): TraitValue | undefined {
  return device?.state[trait];
}

export function isOn(device: DeviceDto | undefined): boolean {
  return traitValue(device, TRAIT_POWER) === true;
}

export function numberOf(device: DeviceDto | undefined, trait: string): number | undefined {
  const value = traitValue(device, trait);

  return typeof value === "number" ? value : undefined;
}

export function brightnessRange(device: DeviceDto | undefined): Range | null {
  return rangeOf(capabilityFor(device, TRAIT_BRIGHTNESS));
}

export function kelvinRange(device: DeviceDto | undefined): Range | null {
  return rangeOf(capabilityFor(device, TRAIT_COLOR_TEMP));
}

export function hasColor(device: DeviceDto | undefined): boolean {
  return supports(device, TRAIT_COLOR);
}

const TRAIT_LABELS: Record<string, string> = {
  [TRAIT_POWER]: "Power",
  [TRAIT_BRIGHTNESS]: "Brightness",
  [TRAIT_COLOR_TEMP]: "Warmth",
  [TRAIT_COLOR]: "Colour",
  fan_speed: "Speed",
  position: "Position",
  target_temperature: "Target",
  temperature: "Temperature",
  humidity: "Humidity",
  battery: "Battery",
  occupancy: "Occupancy",
  power_draw: "Power draw",
  mode: "Mode",
};

export function traitLabel(trait: string): string {
  const known = TRAIT_LABELS[trait];
  if (known) return known;

  return trait
    .split(/[:_]/)
    .filter(Boolean)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}
