import { cn } from "@ava/ui";
import type { CSSProperties, ReactNode } from "react";

import { kelvinToCss } from "@/shared/lib/kelvin";

// First run lights a dark room step by step, and the light is the only colour on screen.
export type RoomLight = {
  /** 0-100. Drives bloom size, opacity and how far the vignette closes in. */
  level: number;
  kelvin: number;
};

export const ROOM_STAGES = {
  asleep: { level: 12, kelvin: 1800 },
  waking: { level: 42, kelvin: 2300 },
  lit: { level: 72, kelvin: 2700 },
  full: { level: 100, kelvin: 2700 },
} as const satisfies Record<string, RoomLight>;

export function Room({
  light,
  children,
  className,
}: {
  light: RoomLight;
  children: ReactNode;
  className?: string;
}) {
  const level = Math.min(Math.max(light.level, 0), 100);

  return (
    <main
      className={cn(
        "relative grid min-h-dvh place-items-center overflow-hidden bg-bg",
        "px-6 py-12",
        className,
      )}
      style={
        {
          "--lit": kelvinToCss(light.kelvin),
          "--level": level,
          // Quadratic so the early steps stay dim: linear opacity reads as a grey wash.
          "--bloom-opacity": `calc(${level} / 100 * ${level} / 100 * 0.62)`,
          // Colour temperature is the slowest thing here: a room warming, not a value changing.
          transition: "background-color 1200ms var(--ease-in-out-quart)",
        } as CSSProperties
      }
    >
      <Bloom level={level} />
      <Vignette level={level} />

      <div className="relative z-raised w-full max-w-[620px]">{children}</div>
    </main>
  );
}

function Bloom({ level }: { level: number }) {
  // Ceiling light: high, wide, and it grows as the room comes up.
  const spread = 620 + level * 7;

  return (
    <div aria-hidden className="pointer-events-none absolute inset-0">
      <div
        className="absolute left-1/2 top-0 animate-bloom-in rounded-full"
        style={{
          width: spread,
          height: spread * 0.78,
          maxWidth: "185vw",
          background: "radial-gradient(closest-side, var(--lit), transparent 74%)",
          animation:
            "bloom-in 1400ms var(--ease-out-expo) both, breathe 7s ease-in-out 1400ms infinite",
          transition:
            "width 1200ms var(--ease-in-out-quart), height 1200ms var(--ease-in-out-quart)",
        }}
      />
    </div>
  );
}

// The darkness closing in from the edges lifts as the room fills, so it reads as a room.
function Vignette({ level }: { level: number }) {
  return (
    <div
      aria-hidden
      className="pointer-events-none absolute inset-0"
      style={{
        background: `radial-gradient(120% 90% at 50% 8%, transparent ${18 + level * 0.42}%, var(--color-bg) 100%)`,
        transition: "background 1200ms var(--ease-in-out-quart)",
      }}
    />
  );
}
