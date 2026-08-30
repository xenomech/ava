import { Materials, coilPath } from "../materials";

import { LIT_SOFT, LIT_SOURCE, type ShapeProps } from "./lit";

/* Fixed geometry, computed once rather than per render. */
const BULB_COIL = coilPath(84, 136, 126, 9, 7);

export function Bulb({ m, className, style }: ShapeProps) {
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
