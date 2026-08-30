import type { CSSProperties } from "react";

import type { useMaterials } from "../materials";

/** What every shape receives: its gradient/filter ids, plus pass-throughs. */
export type ShapeProps = {
  m: ReturnType<typeof useMaterials>;
  className?: string;
  style?: CSSProperties;
};

// How strongly a part glows at --level, shared so equal brightness reads equally.
const litCore = {
  opacity: "calc(var(--level) / 100 * var(--level) / 100 * 0.30)",
} satisfies CSSProperties;
const litSoft = { opacity: "calc(0.05 + var(--level) / 100 * 0.6)" } satisfies CSSProperties;
const litSource = { opacity: "calc(0.14 + var(--level) / 100 * 0.86)" } satisfies CSSProperties;

const TRANSITION = { transition: "opacity 420ms var(--motion-out-soft)" } satisfies CSSProperties;

// Pre-merged: spreading these per node per render allocated hundreds of identical objects.
export const LIT_CORE = { ...litCore, ...TRANSITION } satisfies CSSProperties;
export const LIT_SOFT = { ...litSoft, ...TRANSITION } satisfies CSSProperties;
export const LIT_SOURCE = { ...litSource, ...TRANSITION } satisfies CSSProperties;
export const LIT_FULL = {
  opacity: "calc(var(--level) / 100)",
  ...TRANSITION,
} satisfies CSSProperties;
export const LIT_HALF = {
  opacity: "calc(var(--level) / 100 * .55)",
  ...TRANSITION,
} satisfies CSSProperties;
