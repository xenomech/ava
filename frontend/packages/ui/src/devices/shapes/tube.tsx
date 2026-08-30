import { Materials } from "../materials";

import { LIT_SOFT, LIT_SOURCE, type ShapeProps } from "./lit";

export function Tube({ m, className, style }: ShapeProps) {
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
