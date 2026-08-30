import { Materials } from "../materials";

import { LIT_SOURCE, type ShapeProps } from "./lit";

export function Plug({ m, className, style }: ShapeProps) {
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
