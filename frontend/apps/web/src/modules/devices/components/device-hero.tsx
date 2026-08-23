import { Device, DeviceHalo, cn } from "@ava/ui";
import { emitsLight, isOn, type DeviceDto } from "@ava/contracts";
import { useRef, type CSSProperties, type KeyboardEvent } from "react";

import { useDeviceGesture } from "@/shared/hooks/use-device-gesture";
import { deviceColor, deviceKind } from "../device-view";
import { OnAPlug } from "./on-a-plug";

const KEY_STEP = 5;

export type HeroProps = {
  device: DeviceDto;
  /** Live value while a drag is in flight, so the light tracks the finger. */
  level: number;
  dimmable: boolean;
  offline: boolean;
  onToggle: () => void;
  onDim: (level: number) => void;
  onDimEnd: (level: number) => void;
  onStep: (direction: 1 | -1) => void;
  onOpenSheet: () => void;
};

// The light is the control: the fixture itself is the target of every gesture,
// with no slider laid over it. Everything here is also reachable from the
// keyboard and duplicated as ordinary controls in the sheet, so the gestures
// are the fast path rather than the only path.
export function DeviceHero({
  device,
  level,
  dimmable,
  offline,
  onToggle,
  onDim,
  onDimEnd,
  onStep,
  onOpenSheet,
}: HeroProps) {
  // The hook reports drag distance from where the finger went down, so the level
  // it started from has to survive the gesture. Only that needs a ref: `level`
  // is a prop, and these handlers are rebuilt every render, so they already read
  // the current value without one.
  const from = useRef<number | null>(null);

  const gesture = useDeviceGesture({
    onTap: () => {
      if (!offline) onToggle();
    },
    onHold: onOpenSheet,
    onDim: (delta) => {
      if (offline || !dimmable) return;

      from.current ??= level;
      onDim(clamp(Math.round(from.current + delta)));
    },
    onDimEnd: () => {
      if (from.current === null) return;

      from.current = null;
      onDimEnd(level);
    },
    onSwipe: onStep,
  });

  const color = deviceColor(device);
  const on = isOn(device);
  const lit = emitsLight(deviceKind(device));

  const onKeyDown = (event: KeyboardEvent<HTMLElement>) => {
    const nudge = (direction: 1 | -1) => {
      if (offline || !dimmable) return;

      event.preventDefault();
      onDimEnd(clamp(level + direction * KEY_STEP));
    };

    switch (event.key) {
      case " ":
      case "Enter":
        event.preventDefault();
        if (!offline) onToggle();
        break;
      case "ArrowUp":
        nudge(1);
        break;
      case "ArrowDown":
        nudge(-1);
        break;
      case "ArrowLeft":
        event.preventDefault();
        onStep(-1);
        break;
      case "ArrowRight":
        event.preventDefault();
        onStep(1);
        break;
    }
  };

  return (
    <button
      type="button"
      aria-pressed={on}
      aria-label={`${device.name} — ${reading(device, level, dimmable, offline)}`}
      disabled={offline}
      onKeyDown={onKeyDown}
      {...gesture}
      style={{ ...gesture.style, "--level": level, "--lit": color } as CSSProperties}
      className={cn(
        "group relative grid min-h-0 w-full place-items-center rounded-2xl",
        "focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-fg",
        offline ? "cursor-not-allowed opacity-45" : "cursor-pointer",
      )}
    >
      {lit ? <DeviceHalo className="w-[56%] max-w-[520px]" /> : null}

      {device.kind === "plug" && device.appliance ? (
        <OnAPlug className="absolute right-3 top-3 size-7" />
      ) : null}

      <Device
        kind={deviceKind(device)}
        level={level}
        color={color}
        className={cn(
          "h-[72%] max-h-[420px] transition-transform duration-200 ease-out-soft",
          !offline && "group-active:scale-[0.985]",
        )}
      />
    </button>
  );
}

function reading(device: DeviceDto, level: number, dimmable: boolean, offline: boolean): string {
  if (offline) return "offline";
  if (!isOn(device)) return "off";

  return dimmable ? `${level}% brightness` : "on";
}

function clamp(value: number) {
  return Math.min(Math.max(value, 0), 100);
}
