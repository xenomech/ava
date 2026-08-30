import { Materials } from "../materials";

import { LIT_SOURCE, type ShapeProps } from "./lit";

export function Heater({ m, className, style }: ShapeProps) {
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
