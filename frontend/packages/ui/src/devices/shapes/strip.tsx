import { Materials } from "../materials";

import { LIT_SOFT, LIT_SOURCE, type ShapeProps } from "./lit";

export function Strip({ m, className, style }: ShapeProps) {
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
