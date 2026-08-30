import { Materials } from "../materials";

import { LIT_FULL, LIT_HALF, type ShapeProps } from "./lit";

const FAN_SPOKES = Array.from({ length: 20 }, (_, i) => (i * 360) / 20);

export function Fan({ m, className, style }: ShapeProps) {
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
