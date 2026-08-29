import type { CSSProperties } from "react";

import { cn } from "../lib/utils";
import { Materials, coilPath, useMaterials } from "./materials";

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

type ShapeProps = {
  m: ReturnType<typeof useMaterials>;
  className?: string;
  style?: CSSProperties;
};

const litCore = {
  opacity: "calc(var(--level) / 100 * var(--level) / 100 * 0.30)",
} satisfies CSSProperties;
const litSoft = { opacity: "calc(0.05 + var(--level) / 100 * 0.6)" } satisfies CSSProperties;
const litSource = { opacity: "calc(0.14 + var(--level) / 100 * 0.86)" } satisfies CSSProperties;

const TRANSITION = { transition: "opacity 420ms var(--motion-out-soft)" } satisfies CSSProperties;

/* Pre-merged: the design page renders 27 devices at once, and spreading these
   per node per render allocated hundreds of identical objects. */
const LIT_CORE = { ...litCore, ...TRANSITION } satisfies CSSProperties;
const LIT_SOFT = { ...litSoft, ...TRANSITION } satisfies CSSProperties;
const LIT_SOURCE = { ...litSource, ...TRANSITION } satisfies CSSProperties;
const LIT_FULL = { opacity: "calc(var(--level) / 100)", ...TRANSITION } satisfies CSSProperties;
const LIT_HALF = {
  opacity: "calc(var(--level) / 100 * .55)",
  ...TRANSITION,
} satisfies CSSProperties;

/* Fixed geometry, computed once rather than per render. */
const BULB_COIL = coilPath(84, 136, 126, 9, 7);
const FAN_SPOKES = Array.from({ length: 20 }, (_, i) => (i * 360) / 20);
const HEATER_BOLTS = [0, 1, 2, 3].map((i) => ({
  cx: 130 + 52 * Math.cos(((45 + i * 90) * Math.PI) / 180),
  cy: 212 + 52 * Math.sin(((45 + i * 90) * Math.PI) / 180),
}));

function Bulb({ m, className, style }: ShapeProps) {
  const glass =
    "M110 16c50 0 84 38 84 88 0 27-13 46-25 62-9 12-14 20-14 30v10H65v-10c0-10-5-18-14-30-12-16-25-35-25-62 0-50 34-88 84-88Z";
  const coil = BULB_COIL;

  return (
    <svg viewBox="0 0 220 366" aria-hidden className={className} style={style}>
      <Materials m={m} />

      <path d={glass} fill={`url(#${m.env})`} />

      <ellipse
        cx="110"
        cy="120"
        rx="66"
        ry="70"
        fill="var(--lit)"
        filter={`url(#${m.soft})`}
        style={LIT_SOFT}
      />

      <path d="M104 206v-44h12v44Z" fill="var(--palette-glass)" opacity="0.07" />
      <ellipse cx="110" cy="160" rx="13" ry="8" fill="var(--palette-glass)" opacity="0.06" />

      <g style={LIT_SOURCE}>
        <g stroke="var(--palette-wire)" strokeWidth="1.5" fill="none" opacity="0.85">
          <path d="M100 162V132M120 162v-30" />
          <path d="M92 150l8-6M128 150l-8-6" />
        </g>
        <g fill="none" strokeLinecap="round">
          <path
            d={coil}
            stroke="var(--lit)"
            strokeWidth="9"
            opacity="0.55"
            filter={`url(#${m.glow})`}
          />
          <path d={coil} stroke="var(--lit)" strokeWidth="3.4" />
          <path d={coil} stroke="#fff" strokeWidth="1.3" opacity="0.85" />
          <path d="M84 126l-4 30M136 126l4 30" stroke="var(--lit)" strokeWidth="2.6" />
        </g>
      </g>

      <ellipse
        cx="70"
        cy="74"
        rx="12"
        ry="34"
        fill="#fff"
        opacity="0.26"
        transform="rotate(-24 70 74)"
      />
      <ellipse
        cx="152"
        cy="104"
        rx="5"
        ry="19"
        fill="#fff"
        opacity="0.12"
        transform="rotate(16 152 104)"
      />
      <ellipse cx="110" cy="196" rx="34" ry="7" fill="#fff" style={LIT_SOFT} />

      <path d={glass} fill="none" stroke={`url(#${m.rim})`} strokeWidth="var(--rim-width, 1.6)" />

      <path d="M64 206h92l-4 12H68Z" fill={`url(#${m.dark})`} />

      {[0, 1, 2, 3, 4].map((i) => {
        const y = 220 + i * 13;
        return (
          <g key={i}>
            <path
              d={`M68 ${y + 3} L152 ${y} L152 ${y + 10} L68 ${y + 13} Z`}
              fill={`url(#${m.metal})`}
            />
            <path
              d={`M68 ${y + 3} L152 ${y}`}
              stroke="#fff"
              strokeOpacity="0.4"
              strokeWidth="1"
              fill="none"
            />
            <path
              d={`M68 ${y + 13} L152 ${y + 10}`}
              stroke="#000"
              strokeOpacity="0.55"
              strokeWidth="1.6"
              fill="none"
            />
          </g>
        );
      })}

      <path d="M78 288h64v14a14 14 0 0 1-14 14H92a14 14 0 0 1-14-14Z" fill={`url(#${m.dark})`} />
      <ellipse cx="110" cy="320" rx="15" ry="9" fill={`url(#${m.metal})`} />
    </svg>
  );
}

function Strip({ m, className, style }: ShapeProps) {
  return (
    <svg viewBox="0 0 320 366" aria-hidden className={className} style={style}>
      <Materials m={m} />

      <g transform="rotate(-6 160 184)">
        <rect x="30" y="160" width="260" height="38" rx="19" fill={`url(#${m.dark})`} />
        <rect x="30" y="160" width="260" height="14" rx="7" fill="#fff" opacity="0.13" />
        <rect x="30" y="188" width="260" height="10" rx="5" fill="#000" opacity="0.3" />

        <rect x="42" y="170" width="236" height="18" rx="9" fill="var(--lit)" style={LIT_SOFT} />

        <g style={LIT_SOURCE}>
          {[0, 1, 2, 3, 4, 5].map((i) => {
            const x = 68 + i * 37;
            return (
              <g key={i}>
                <circle
                  cx={x}
                  cy="179"
                  r="9"
                  fill="var(--lit)"
                  opacity="0.6"
                  filter={`url(#${m.glow})`}
                />
                <rect x={x - 5} y="174" width="10" height="10" rx="2" fill="var(--lit)" />
                <rect
                  x={x - 2.5}
                  y="176.5"
                  width="5"
                  height="5"
                  rx="1"
                  fill="#fff"
                  opacity="0.75"
                />
              </g>
            );
          })}
        </g>

        <rect x="26" y="166" width="10" height="26" rx="4" fill={`url(#${m.metal})`} />
      </g>
    </svg>
  );
}

function Lamp({ m, className, style }: ShapeProps) {
  return (
    <svg viewBox="0 0 260 366" aria-hidden className={className} style={style}>
      <Materials m={m} />
      <path
        d="M80 152h100l46 146H34Z"
        fill="var(--lit)"
        filter={`url(#${m.soft})`}
        style={LIT_CORE}
      />

      <path d="M78 78h104l22 72H56Z" fill={`url(#${m.dark})`} />
      <g opacity="0.13" stroke="#fff" strokeWidth="0.8">
        {[0, 1, 2, 3, 4, 5, 6].map((i) => (
          <path key={i} d={`M${80 + i * 14} 78 L${72 + i * 16} 150`} />
        ))}
      </g>
      <path d="M78 78h104l3 10H75Z" fill="#fff" opacity="0.15" />

      <path d="M60 142h140l-4 8H64Z" fill="var(--lit)" style={LIT_SOFT} />
      <g style={LIT_SOURCE}>
        <ellipse cx="130" cy="150" rx="72" ry="10" fill="var(--lit)" filter={`url(#${m.glow})`} />
        <ellipse cx="130" cy="150" rx="52" ry="6" fill="#fff" opacity="0.4" />
      </g>

      <rect x="123" y="150" width="14" height="150" fill={`url(#${m.metal})`} />
      <rect x="123" y="150" width="4" height="150" fill="#fff" opacity="0.16" />
      <ellipse cx="130" cy="304" rx="62" ry="15" fill={`url(#${m.dark})`} />
      <ellipse cx="130" cy="299" rx="62" ry="15" fill={`url(#${m.metal})`} opacity="0.9" />
    </svg>
  );
}

function Plug({ m, className, style }: ShapeProps) {
  return (
    <svg viewBox="0 0 260 366" aria-hidden className={className} style={style}>
      <Materials m={m} />

      <rect x="104" y="72" width="14" height="42" rx="4" fill={`url(#${m.metal})`} />
      <rect x="142" y="72" width="14" height="42" rx="4" fill={`url(#${m.metal})`} />

      <rect x="58" y="112" width="144" height="156" rx="34" fill={`url(#${m.dark})`} />
      <rect x="58" y="112" width="144" height="54" rx="34" fill="#fff" opacity="0.08" />
      <rect
        x="58"
        y="112"
        width="144"
        height="156"
        rx="34"
        fill="none"
        stroke={`url(#${m.rim})`}
        strokeWidth="var(--rim-width, 1.4)"
      />

      <g style={LIT_SOURCE}>
        <circle cx="130" cy="168" r="17" fill="var(--lit)" filter={`url(#${m.glow})`} />
        <circle cx="130" cy="168" r="9" fill="#fff" opacity="0.55" />
      </g>

      <rect x="96" y="208" width="68" height="12" rx="6" fill="#08090b" />
      <rect x="112" y="232" width="36" height="10" rx="5" fill="#08090b" opacity="0.8" />
    </svg>
  );
}

function Sensor({ m, className, style }: ShapeProps) {
  return (
    <svg viewBox="0 0 260 366" aria-hidden className={className} style={style}>
      <Materials m={m} />

      <circle cx="130" cy="186" r="78" fill={`url(#${m.dark})`} />
      <path d="M52 186a78 78 0 0 1 156 0Z" fill="#fff" opacity="0.07" />
      <circle
        cx="130"
        cy="186"
        r="78"
        fill="none"
        stroke={`url(#${m.rim})`}
        strokeWidth="var(--rim-width, 1.4)"
      />

      <circle cx="130" cy="186" r="50" fill={`url(#${m.env})`} />
      <g fill="none" stroke="#fff" strokeOpacity="0.18" strokeWidth="0.9">
        {[42, 34, 26, 18].map((r) => (
          <circle key={r} cx="130" cy="186" r={r} />
        ))}
      </g>

      <circle
        cx="130"
        cy="186"
        r="12"
        fill="var(--lit)"
        filter={`url(#${m.glow})`}
        style={LIT_SOURCE}
      />

      <rect x="118" y="264" width="24" height="16" rx="4" fill={`url(#${m.metal})`} />
      <rect x="106" y="278" width="48" height="8" rx="4" fill={`url(#${m.dark})`} />
    </svg>
  );
}

function Tube({ m, className, style }: ShapeProps) {
  const emit = `tube-emit-${m.id}`;
  const spill = `tube-spill-${m.id}`;

  return (
    <svg viewBox="0 0 380 366" aria-hidden className={className} style={style}>
      <Materials m={m} />

      <defs>
        <linearGradient id={emit} x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor="var(--lit)" stopOpacity="0.35" />
          <stop offset="45%" stopColor="var(--lit)" stopOpacity="1" />
          <stop offset="100%" stopColor="var(--lit)" stopOpacity="0.55" />
        </linearGradient>

        <radialGradient id={spill} cx="50%" cy="0%" r="100%">
          <stop offset="0%" stopColor="var(--lit)" stopOpacity="0.45" />
          <stop offset="100%" stopColor="var(--lit)" stopOpacity="0" />
        </radialGradient>
      </defs>

      <rect x="16" y="156" width="348" height="46" rx="15" fill={`url(#${m.env})`} />
      <rect x="16" y="156" width="348" height="15" rx="7" fill="#fff" opacity="0.11" />
      <rect x="16" y="190" width="348" height="12" rx="6" fill="#000" opacity="0.28" />

      <rect
        x="30"
        y="168"
        width="320"
        height="22"
        rx="11"
        fill={`url(#${emit})`}
        style={LIT_SOFT}
      />

      <rect
        x="30"
        y="168"
        width="320"
        height="22"
        rx="11"
        fill="var(--lit)"
        filter={`url(#${m.glow})`}
        style={LIT_SOURCE}
      />

      <rect
        x="16"
        y="156"
        width="348"
        height="46"
        rx="15"
        fill="none"
        stroke={`url(#${m.rim})`}
        strokeWidth="var(--rim-width, 1.6)"
      />

      <ellipse cx="190" cy="212" rx="176" ry="58" fill={`url(#${spill})`} style={LIT_SOFT} />
    </svg>
  );
}

function Fan({ m, className, style }: ShapeProps) {
  const spokes = FAN_SPOKES;
  const blade = "M130 148 C168 128 176 92 158 58 C140 74 116 92 130 148 Z";

  return (
    <svg viewBox="0 0 260 366" aria-hidden className={className} style={style}>
      <Materials m={m} />

      <ellipse cx="130" cy="338" rx="70" ry="9" fill="#000" opacity=".5" />

      <rect x="121" y="230" width="18" height="88" rx="4" fill={`url(#${m.metal})`} />
      <rect x="126" y="230" width="3" height="88" fill="#fff" opacity=".2" />
      <ellipse cx="130" cy="318" rx="46" ry="11" fill={`url(#${m.metal})`} />
      <ellipse cx="130" cy="314" rx="46" ry="11" fill={`url(#${m.dark})`} />

      <circle cx="130" cy="148" r="86" fill="#0a0b0d" />

      <g style={LIT_HALF}>
        <circle
          cx="130"
          cy="148"
          r="72"
          fill="none"
          stroke="#8a8a94"
          strokeOpacity=".3"
          strokeWidth="30"
        />
      </g>

      <g>
        {[0, 120, 240].map((a) => (
          <g key={a} transform={`rotate(${a} 130 148)`}>
            <path d={blade} fill="#2b3038" />
            <path d={blade} fill={`url(#${m.env})`} opacity=".7" />
            <path
              d="M130 148 C156 130 164 102 152 74"
              fill="none"
              stroke="#fff"
              strokeOpacity=".16"
              strokeWidth="2"
            />
          </g>
        ))}
      </g>

      <circle cx="130" cy="148" r="20" fill={`url(#${m.metal})`} />
      <circle
        cx="130"
        cy="148"
        r="20"
        fill="none"
        stroke="#000"
        strokeOpacity=".45"
        strokeWidth="1.6"
      />
      <circle cx="124" cy="142" r="5" fill="#fff" opacity=".28" />

      <g stroke="#aeb5bf" fill="none" strokeLinecap="round" opacity=".55">
        {spokes.map((a) => (
          <path key={a} d="M130 66 V148" strokeWidth="1.6" transform={`rotate(${a} 130 148)`} />
        ))}
        <circle cx="130" cy="148" r="58" strokeWidth="1.8" />
        <circle cx="130" cy="148" r="32" strokeWidth="1.8" />
      </g>

      <circle cx="130" cy="148" r="86" fill="none" stroke={`url(#${m.metal})`} strokeWidth="7" />
      <circle
        cx="130"
        cy="148"
        r="86"
        fill="none"
        stroke="#000"
        strokeOpacity=".35"
        strokeWidth="1.4"
      />

      <circle cx="130" cy="322" r="4" fill="#4ade80" style={LIT_FULL} />
    </svg>
  );
}

function Heater({ m, className, style }: ShapeProps) {
  return (
    <svg viewBox="0 0 260 366" aria-hidden className={className} style={style}>
      <Materials m={m} />

      <rect x="42" y="58" width="176" height="238" rx="18" fill={`url(#${m.env})`} />
      <rect x="42" y="58" width="176" height="16" rx="8" fill="#fff" opacity=".1" />

      <g style={LIT_SOURCE}>
        {[0, 1, 2, 3, 4].map((i) => {
          const y = 96 + i * 40;
          return (
            <g key={i}>
              <rect
                x="66"
                y={y}
                width="128"
                height="12"
                rx="6"
                fill="var(--lit)"
                filter={`url(#${m.glow})`}
                opacity=".7"
              />
              <rect x="66" y={y} width="128" height="12" rx="6" fill="var(--lit)" />
              <rect x="72" y={y + 3} width="116" height="3" rx="1.5" fill="#fff" opacity=".5" />
            </g>
          );
        })}
      </g>

      <g stroke={`url(#${m.metal})`} strokeWidth="4" fill="none">
        {[0, 1, 2, 3, 4, 5].map((i) => (
          <path key={i} d={`M${58 + i * 29} 58 V296`} strokeOpacity=".5" />
        ))}
      </g>

      <rect
        x="42"
        y="58"
        width="176"
        height="238"
        rx="18"
        fill="none"
        stroke={`url(#${m.rim})`}
        strokeWidth="1.6"
      />
      <path d="M58 296h144l10 24H48Z" fill={`url(#${m.dark})`} />
    </svg>
  );
}

function Speaker({ m, className, style }: ShapeProps) {
  return (
    <svg viewBox="0 0 260 366" aria-hidden className={className} style={style}>
      <Materials m={m} />

      <ellipse cx="130" cy="322" rx="76" ry="9" fill="#000" opacity=".45" />

      <rect x="52" y="46" width="156" height="272" rx="8" fill={`url(#${m.dark})`} />
      <rect x="52" y="46" width="156" height="272" rx="8" fill={`url(#${m.env})`} opacity=".5" />
      <rect x="52" y="46" width="156" height="7" rx="3" fill="#fff" opacity=".12" />
      <rect x="52" y="311" width="156" height="7" rx="3" fill="#000" opacity=".4" />

      <circle cx="130" cy="212" r="62" fill="#0d0f12" />
      <circle
        cx="130"
        cy="212"
        r="62"
        fill="none"
        stroke="#fff"
        strokeOpacity=".07"
        strokeWidth="2"
      />
      <circle cx="130" cy="212" r="52" fill={`url(#${m.env})`} opacity=".5" />
      <circle
        cx="130"
        cy="212"
        r="52"
        fill="none"
        stroke="#000"
        strokeOpacity=".55"
        strokeWidth="7"
      />
      <circle cx="130" cy="212" r="38" fill="#15181c" />
      <circle
        cx="130"
        cy="212"
        r="38"
        fill="none"
        stroke="#fff"
        strokeOpacity=".05"
        strokeWidth="1.5"
      />
      <circle cx="130" cy="212" r="17" fill={`url(#${m.metal})`} />
      <ellipse
        cx="123"
        cy="204"
        rx="6"
        ry="9"
        fill="#fff"
        opacity=".22"
        transform="rotate(-24 123 204)"
      />

      {HEATER_BOLTS.map((bolt) => (
        <circle key={bolt.cx} cx={bolt.cx} cy={bolt.cy} r="3.4" fill={`url(#${m.metal})`} />
      ))}

      <circle cx="130" cy="104" r="26" fill="#0d0f12" />
      <circle
        cx="130"
        cy="104"
        r="26"
        fill="none"
        stroke="#fff"
        strokeOpacity=".07"
        strokeWidth="2"
      />
      <circle cx="130" cy="104" r="13" fill={`url(#${m.metal})`} />
      <ellipse
        cx="125"
        cy="99"
        rx="4"
        ry="6"
        fill="#fff"
        opacity=".24"
        transform="rotate(-24 125 99)"
      />

      <circle cx="196" cy="104" r="7" fill="#08090b" />
      <circle cx="196" cy="104" r="7" fill="none" stroke="#fff" strokeOpacity=".06" />

      <rect
        x="52"
        y="46"
        width="156"
        height="272"
        rx="8"
        fill="none"
        stroke={`url(#${m.rim})`}
        strokeWidth="1.6"
      />

      <circle cx="64" cy="304" r="3.6" fill="#4ade80" style={LIT_FULL} />
    </svg>
  );
}

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

const HALO_STYLE = {
  background: "radial-gradient(circle, var(--lit) 0%, transparent 68%)",
  opacity: "calc(var(--level) / 100 * var(--level) / 100 * 0.2)",
  transition: "opacity 500ms var(--motion-out-soft)",
} satisfies CSSProperties;
