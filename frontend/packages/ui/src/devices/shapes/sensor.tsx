import { Materials } from "../materials";

import { LIT_SOURCE, type ShapeProps } from "./lit";

export function Sensor({ m, className, style }: ShapeProps) {
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
