import type { DeviceDto } from "./device.dto";

export const DEVICE_FORMS = ["bulb", "tube", "strip", "lamp", "plug", "sensor"] as const;

export type DeviceForm = (typeof DEVICE_FORMS)[number];

export type DeviceProfile = {
  kind: DeviceForm;
  form: string;
};

const MODELS: Record<string, DeviceProfile> = {
  ESP25_SHTW_01: { kind: "tube", form: "Batten" },
};

const KINDS: Record<string, DeviceProfile> = {
  bulb: { kind: "bulb", form: "Bulb" },
  tube: { kind: "tube", form: "Batten" },
  strip: { kind: "strip", form: "Light strip" },
  lamp: { kind: "lamp", form: "Lamp" },
  plug: { kind: "plug", form: "Plug" },
  sensor: { kind: "sensor", form: "Sensor" },
};

const UNKNOWN: DeviceProfile = { kind: "bulb", form: "Light" };

export function deviceProfile(device: DeviceDto): DeviceProfile {
  const model = device.state.model;

  if (model && MODELS[model]) return MODELS[model];

  return KINDS[device.kind] ?? UNKNOWN;
}

export type Range = { min: number; max: number };

export type DeviceControls = {
  brightness: Range | null;
  kelvin: Range | null;
  color: boolean;
};

const NONE: DeviceControls = { brightness: null, kelvin: null, color: false };

export function deviceControls(device: DeviceDto | undefined): DeviceControls {
  if (!device) return NONE;

  const { capabilities, limits } = device.state;

  return {
    brightness: capabilities.includes("brightness")
      ? { min: limits?.brightness_min ?? 0, max: limits?.brightness_max ?? 100 }
      : null,
    kelvin:
      capabilities.includes("color_temp") &&
      limits?.kelvin_min !== undefined &&
      limits?.kelvin_max !== undefined
        ? { min: limits.kelvin_min, max: limits.kelvin_max }
        : null,
    color: capabilities.includes("color"),
  };
}
