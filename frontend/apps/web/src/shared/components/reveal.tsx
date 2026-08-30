import { cn } from "@ava/ui";
import type { CSSProperties, ReactNode } from "react";

// One line of the cascade. Nothing on a first-run screen should arrive at the
// same moment as the thing above it — the offset is most of what separates a
// considered entrance from a page that simply appeared.
//
// Under prefers-reduced-motion the global stylesheet collapses both duration
// and delay, so everything lands immediately and nothing is left invisible.
const STEP_MS = 90;

export function Reveal({
  at = 0,
  children,
  className,
}: {
  /** Position in the cascade, not milliseconds. */
  at?: number;
  children: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={cn("animate-rise", className)}
      style={{ animationDelay: `${at * STEP_MS}ms` } as CSSProperties}
    >
      {children}
    </div>
  );
}
