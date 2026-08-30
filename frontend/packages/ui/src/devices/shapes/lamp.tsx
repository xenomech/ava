import { Materials } from "../materials";

import { LIT_CORE, LIT_SOFT, LIT_SOURCE, type ShapeProps } from "./lit";

export function Lamp({ m, className, style }: ShapeProps) {
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
