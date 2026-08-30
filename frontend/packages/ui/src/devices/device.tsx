import type { CSSProperties } from "react";

import { cn } from "../lib/utils";
import { Bulb } from "./shapes/bulb";
import { Fan } from "./shapes/fan";
import { Heater } from "./shapes/heater";
import { Lamp } from "./shapes/lamp";
import { Plug } from "./shapes/plug";
import { Sensor } from "./shapes/sensor";
import { Speaker } from "./shapes/speaker";
import { Strip } from "./shapes/strip";
import { Tube } from "./shapes/tube";
import { useMaterials } from "./materials";
import type { ShapeProps } from "./shapes/lit";

export type DeviceKind =
  | "bulb"
  | "tube"
  | "strip"
  | "lamp"
  | "plug"
  | "sensor"
  | "fan"
  | "heater"
  | "speaker";

export type DeviceProps = {
  kind: DeviceKind;
  level?: number;
  color?: string;
  className?: string;
  style?: CSSProperties;
};

const SHAPES: Record<DeviceKind, (props: ShapeProps) => React.ReactElement> = {
  bulb: Bulb,
  tube: Tube,
  strip: Strip,
  lamp: Lamp,
  plug: Plug,
  sensor: Sensor,
  fan: Fan,
  heater: Heater,
  speaker: Speaker,
};

/** Dispatcher for a runtime `kind`. A caller that knows its fixture statically
    should import that shape directly and skip the other eight. */
export function Device({ kind, level = 0, color = "#ffb463", className, style }: DeviceProps) {
  const m = useMaterials();
  const Shape = SHAPES[kind];

  return (
    <Shape
      m={m}
      className={cn("h-full w-auto max-h-full max-w-full", className)}
      style={{ "--level": level, "--lit": color, ...style } as CSSProperties}
    />
  );
}

const HALO_STYLE = {
  background: "radial-gradient(circle, var(--lit) 0%, transparent 68%)",
  opacity: "calc(var(--level) / 100 * var(--level) / 100 * 0.2)",
  transition: "opacity 500ms var(--motion-out-soft)",
} satisfies CSSProperties;

export function DeviceHalo({ className }: { className?: string }) {
  return (
    <span
      aria-hidden
      className={cn(
        "pointer-events-none absolute aspect-square rounded-full blur-[52px]",
        className,
      )}
      style={HALO_STYLE}
    />
  );
}
