import { Materials } from "../materials";

import { LIT_FULL, type ShapeProps } from "./lit";

const SPEAKER_BOLTS = [0, 1, 2, 3].map((i) => ({
  cx: 130 + 52 * Math.cos(((45 + i * 90) * Math.PI) / 180),
  cy: 212 + 52 * Math.sin(((45 + i * 90) * Math.PI) / 180),
}));

export function Speaker({ m, className, style }: ShapeProps) {
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

      {SPEAKER_BOLTS.map((bolt, at) => (
        <circle key={at} cx={bolt.cx} cy={bolt.cy} r="3.4" fill={`url(#${m.metal})`} />
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
