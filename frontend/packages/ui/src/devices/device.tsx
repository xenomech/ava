import type { CSSProperties } from "react";

import { cn } from "../lib/utils";
import { Materials, coilPath, useMaterials } from "./materials";

export type DeviceKind = "bulb" | "strip" | "lamp" | "plug" | "sensor";

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

const litCore = { opacity: "calc(var(--level) / 100 * 0.9)" } satisfies CSSProperties;
const litSoft = { opacity: "calc(0.05 + var(--level) / 100 * 0.6)" } satisfies CSSProperties;
const litSource = { opacity: "calc(0.14 + var(--level) / 100 * 0.86)" } satisfies CSSProperties;

const TRANSITION = { transition: "opacity 420ms var(--motion-out-soft)" } satisfies CSSProperties;

function Bulb({ m, className, style }: ShapeProps) {
  const glass =
    "M110 16c50 0 84 38 84 88 0 27-13 46-25 62-9 12-14 20-14 30v10H65v-10c0-10-5-18-14-30-12-16-25-35-25-62 0-50 34-88 84-88Z";
  const coil = coilPath(84, 136, 126, 9, 7);

  return (
    <svg viewBox="0 0 220 366" aria-hidden className={className} style={style}>
      <Materials m={m} />

      <ellipse
        cx="110"
        cy="124"
        rx="92"
        ry="96"
        fill={`url(#${m.core})`}
        filter={`url(#${m.soft})`}
        style={{ ...litCore, ...TRANSITION }}
      />

      <path d={glass} fill={`url(#${m.env})`} />

      <ellipse
        cx="110"
        cy="120"
        rx="66"
        ry="70"
        fill="var(--lit)"
        filter={`url(#${m.soft})`}
        style={{ ...litSoft, ...TRANSITION }}
      />

      <path d="M104 206v-44h12v44Z" fill="var(--palette-glass)" opacity="0.07" />
      <ellipse cx="110" cy="160" rx="13" ry="8" fill="var(--palette-glass)" opacity="0.06" />

      <g style={{ ...litSource, ...TRANSITION }}>
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
      <ellipse cx="110" cy="196" rx="34" ry="7" fill="#fff" style={{ ...litSoft, ...TRANSITION }} />

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
      <ellipse
        cx="160"
        cy="190"
        rx="146"
        ry="58"
        fill={`url(#${m.core})`}
        filter={`url(#${m.soft})`}
        style={{ ...litCore, ...TRANSITION }}
      />

      <g transform="rotate(-6 160 184)">
        <rect x="30" y="160" width="260" height="38" rx="19" fill={`url(#${m.dark})`} />
        <rect x="30" y="160" width="260" height="14" rx="7" fill="#fff" opacity="0.13" />
        <rect x="30" y="188" width="260" height="10" rx="5" fill="#000" opacity="0.3" />

        <rect
          x="42"
          y="170"
          width="236"
          height="18"
          rx="9"
          fill="var(--lit)"
          style={{ ...litSoft, ...TRANSITION }}
        />

        <g style={{ ...litSource, ...TRANSITION }}>
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
        d="M76 152h108l66 158H10Z"
        fill="var(--lit)"
        filter={`url(#${m.soft})`}
        style={{ ...litCore, ...TRANSITION }}
      />

      <path d="M78 78h104l22 72H56Z" fill={`url(#${m.dark})`} />
      <g opacity="0.13" stroke="#fff" strokeWidth="0.8">
        {[0, 1, 2, 3, 4, 5, 6].map((i) => (
          <path key={i} d={`M${80 + i * 14} 78 L${72 + i * 16} 150`} />
        ))}
      </g>
      <path d="M78 78h104l3 10H75Z" fill="#fff" opacity="0.15" />

      <path d="M60 142h140l-4 8H64Z" fill="var(--lit)" style={{ ...litSoft, ...TRANSITION }} />
      <g style={{ ...litSource, ...TRANSITION }}>
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
      <ellipse
        cx="130"
        cy="200"
        rx="98"
        ry="82"
        fill={`url(#${m.core})`}
        filter={`url(#${m.soft})`}
        style={{ ...litCore, ...TRANSITION }}
      />

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

      <g style={{ ...litSource, ...TRANSITION }}>
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
      <ellipse
        cx="130"
        cy="186"
        rx="94"
        ry="82"
        fill={`url(#${m.core})`}
        filter={`url(#${m.soft})`}
        style={{ ...litCore, ...TRANSITION }}
      />

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
        style={{ ...litSource, ...TRANSITION }}
      />

      <rect x="118" y="264" width="24" height="16" rx="4" fill={`url(#${m.metal})`} />
      <rect x="106" y="278" width="48" height="8" rx="4" fill={`url(#${m.dark})`} />
    </svg>
  );
}

const SHAPES: Record<DeviceKind, (props: ShapeProps) => React.ReactElement> = {
  bulb: Bulb,
  strip: Strip,
  lamp: Lamp,
  plug: Plug,
  sensor: Sensor,
};

export function DeviceHalo({ className }: { className?: string }) {
  return (
    <span
      aria-hidden
      className={cn(
        "pointer-events-none absolute aspect-square rounded-full blur-[46px]",
        className,
      )}
      style={{
        background: "radial-gradient(circle, var(--lit) 0%, transparent 68%)",
        opacity: "calc(var(--level) / 100 * 0.22)",
        transition: "opacity 500ms var(--motion-out-soft)",
      }}
    />
  );
}
