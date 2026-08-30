import { useId, useMemo } from "react";

const COLONS = /:/g;

export function useMaterials() {
  const raw = useId();

  // One stable object per instance, or a memoised shape could never bail out.
  return useMemo(() => {
    const id = raw.replace(COLONS, "");

    return {
      id,
      core: `core-${id}`,
      env: `env-${id}`,
      rim: `rim-${id}`,
      metal: `metal-${id}`,
      dark: `dark-${id}`,
      soft: `soft-${id}`,
      glow: `glow-${id}`,
    };
  }, [raw]);
}

export function Materials({ m }: { m: ReturnType<typeof useMaterials> }) {
  return (
    <defs>
      <radialGradient id={m.core} cx="50%" cy="46%" r="52%">
        <stop offset="0%" stopColor="var(--lit)" stopOpacity="1" />
        <stop offset="34%" stopColor="var(--lit)" stopOpacity="0.7" />
        <stop offset="100%" stopColor="var(--lit)" stopOpacity="0" />
      </radialGradient>

      <radialGradient id={m.env} cx="38%" cy="32%" r="72%">
        <stop offset="0%" stopColor="var(--palette-glass)" stopOpacity="var(--glass-1, 0.2)" />
        <stop offset="46%" stopColor="var(--palette-glass)" stopOpacity="var(--glass-2, 0.045)" />
        <stop offset="82%" stopColor="var(--palette-glass)" stopOpacity="var(--glass-3, 0.02)" />
        <stop offset="100%" stopColor="var(--palette-glass)" stopOpacity="var(--glass-4, 0.13)" />
      </radialGradient>

      <linearGradient id={m.rim} x1="0" y1="0" x2="1" y2="1">
        <stop offset="0%" stopColor="var(--palette-glass)" stopOpacity="var(--rim-1, 0.5)" />
        <stop offset="40%" stopColor="var(--palette-glass)" stopOpacity="var(--rim-2, 0.12)" />
        <stop offset="100%" stopColor="var(--palette-glass)" stopOpacity="var(--rim-3, 0.42)" />
      </linearGradient>

      <linearGradient id={m.metal} x1="0" y1="0" x2="1" y2="0">
        <stop offset="0%" stopColor="#2a2e34" />
        <stop offset="9%" stopColor="#6b727d" />
        <stop offset="21%" stopColor="#ccd3dd" />
        <stop offset="33%" stopColor="#7d848f" />
        <stop offset="47%" stopColor="#eef2f8" />
        <stop offset="61%" stopColor="#79808b" />
        <stop offset="76%" stopColor="#bcc3cd" />
        <stop offset="89%" stopColor="#565d66" />
        <stop offset="100%" stopColor="#22262b" />
      </linearGradient>

      <linearGradient id={m.dark} x1="0" y1="0" x2="1" y2="0">
        <stop offset="0%" stopColor="#101216" />
        <stop offset="42%" stopColor="#2e343c" />
        <stop offset="70%" stopColor="#1c2026" />
        <stop offset="100%" stopColor="#0d0f12" />
      </linearGradient>

      <filter id={m.soft} x="-90%" y="-90%" width="280%" height="280%">
        <feGaussianBlur stdDeviation="11" />
      </filter>
      <filter id={m.glow} x="-70%" y="-70%" width="240%" height="240%">
        <feGaussianBlur stdDeviation="5" />
      </filter>
    </defs>
  );
}

export function coilPath(x1: number, x2: number, y: number, peaks: number, amp: number) {
  const step = (x2 - x1) / (peaks * 2);
  let d = `M${x1} ${y}`;
  for (let i = 0; i < peaks * 2; i++) {
    d += ` L${(x1 + step * (i + 1)).toFixed(1)} ${(y + (i % 2 ? amp : -amp)).toFixed(1)}`;
  }
  return d;
}
