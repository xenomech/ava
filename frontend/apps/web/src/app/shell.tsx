import { cn } from "@ava/ui";
import { Outlet } from "@tanstack/react-router";
import type { ReactNode } from "react";

import { useMediaQuery } from "@/shared/hooks/use-media-query";
import { MenuHandle } from "./menu-handle";

/** Where the rail replaces the pull-down. Matches the `md:` classes below. */
const RAIL = "(min-width: 768px)";

/** Pure layout with no top bar: the composition root decides what fills the slots. */
export function Shell({ nav, palette }: { nav: ReactNode; palette?: ReactNode }) {
  // Mounted per form factor: a display:none rail still runs every query and re-renders.
  const rail = useMediaQuery(RAIL);

  return (
    // pl/pr-safe only: the surfaces that touch the top and bottom edges handle those themselves.
    <div
      className={cn(
        "grid h-dvh bg-surface pl-safe pr-safe",
        "md:grid-cols-[208px_minmax(0,1fr)] xl:grid-cols-[262px_minmax(0,1fr)]",
      )}
    >
      {rail ? (
        <aside className="hidden min-h-0 bg-surface px-4 pb-8 pt-8 md:block xl:px-6 xl:pb-10 xl:pt-10">
          {nav}
        </aside>
      ) : null}

      <main
        className={cn(
          "relative m-2 min-h-0 min-w-0 overflow-y-auto rounded-2xl bg-bg",
          "mb-[max(0.5rem,env(safe-area-inset-bottom))] mt-[max(0.5rem,env(safe-area-inset-top))]",
        )}
      >
        {/* The pull-down renders the same node, so the phone can never drift from the desktop. */}
        {rail ? null : <MenuHandle>{nav}</MenuHandle>}

        <Outlet />
      </main>

      {palette}
    </div>
  );
}
