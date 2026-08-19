import type { DeviceKind } from "@ava/ui";
import { create } from "zustand";
import { useShallow } from "zustand/react/shallow";

import { kelvinToCss } from "@/shared/lib/kelvin";

export type Device = {
  id: string;
  name: string;
  room: string;
  kind: DeviceKind;
  level: number;
  color: string;
  kelvin: number | null;
  watts: number;
  lastLevel: number;
};

type DeviceState = {
  devices: Device[];
  focused: string;

  focus: (id: string) => void;
  step: (direction: 1 | -1) => void;
  setLevel: (id: string, level: number) => void;
  toggle: (id: string) => void;
  setColor: (id: string, color: string) => void;
  setWhite: (id: string, kelvin: number) => void;
};

const clamp = (v: number, lo = 0, hi = 100) => Math.min(Math.max(v, lo), hi);

const SEED: Device[] = [
  {
    id: "ceiling",
    name: "Ceiling",
    room: "Living room",
    kind: "bulb",
    level: 72,
    color: kelvinToCss(2700),
    kelvin: 2700,
    watts: 7.4,
    lastLevel: 72,
  },
  {
    id: "playbar",
    name: "Play bar",
    room: "Living room",
    kind: "strip",
    level: 40,
    color: kelvinToCss(2200),
    kelvin: 2200,
    watts: 11.2,
    lastLevel: 40,
  },
  {
    id: "reading",
    name: "Reading lamp",
    room: "Bedroom",
    kind: "lamp",
    level: 0,
    color: kelvinToCss(3600),
    kelvin: 3600,
    watts: 9,
    lastLevel: 65,
  },
  {
    id: "coffee",
    name: "Coffee machine",
    room: "Kitchen",
    kind: "plug",
    level: 100,
    color: kelvinToCss(3000),
    kelvin: 3000,
    watts: 0.4,
    lastLevel: 100,
  },
  {
    id: "hall",
    name: "Hall sensor",
    room: "Hallway",
    kind: "sensor",
    level: 30,
    color: kelvinToCss(5200),
    kelvin: 5200,
    watts: 0.2,
    lastLevel: 30,
  },
];

export const useDeviceStore = create<DeviceState>()((set, get) => ({
  devices: SEED,
  focused: SEED[0]!.id,

  focus: (id) => set({ focused: id }),

  step: (direction) => {
    const { devices, focused } = get();
    const at = devices.findIndex((d) => d.id === focused);
    const next = devices[clamp(at + direction, 0, devices.length - 1)];
    if (next) set({ focused: next.id });
  },

  setLevel: (id, level) =>
    set((s) => ({
      devices: s.devices.map((d) =>
        d.id === id
          ? {
              ...d,
              level: Math.round(clamp(level)),
              lastLevel: level > 0 ? Math.round(clamp(level)) : d.lastLevel,
            }
          : d,
      ),
    })),

  toggle: (id) =>
    set((s) => ({
      devices: s.devices.map((d) =>
        d.id === id ? { ...d, level: d.level > 0 ? 0 : d.lastLevel || 65 } : d,
      ),
    })),

  setColor: (id, color) =>
    set((s) => ({
      devices: s.devices.map((d) => (d.id === id ? { ...d, color, kelvin: null } : d)),
    })),

  setWhite: (id, kelvin) =>
    set((s) => ({
      devices: s.devices.map((d) =>
        d.id === id ? { ...d, kelvin, color: kelvinToCss(kelvin) } : d,
      ),
    })),
}));

export function useFocusedDevice() {
  return useDeviceStore(
    useShallow((s) => ({
      device: s.devices.find((d) => d.id === s.focused) ?? s.devices[0]!,
      setLevel: s.setLevel,
      toggle: s.toggle,
      setColor: s.setColor,
      setWhite: s.setWhite,
      step: s.step,
    })),
  );
}

export function useDevices() {
  return useDeviceStore(
    useShallow((s) => ({ devices: s.devices, focused: s.focused, focus: s.focus })),
  );
}
