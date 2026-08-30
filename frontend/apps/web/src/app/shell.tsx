import { cn } from "@ava/ui";
import { Outlet } from "@tanstack/react-router";
import type { ReactNode } from "react";

import { useMediaQuery } from "@/shared/hooks/use-media-query";
import { MenuHandle } from "./menu-handle";

/** Where the rail replaces the pull-down. Matches the `md:` classes below. */
const RAIL = "(min-width: 768px)";

/**
 * No top bar. The window is the app: navigation down the side, the light
 * filling everything else.
 *
 * Pure layout: what fills the nav and palette slots is the composition root's
 * business (see app.tsx) — the shell knows nothing about any module.
 *
 * On a phone the sidebar becomes a panel pulled down from the top edge, so the
 * only permanent chrome is a handle the width of two fingers.
 */
export function Shell({ nav, palette }: { nav: ReactNode; palette?: ReactNode }) {
  /* Mounted per form factor rather than hidden with CSS: a display:none rail
     still runs every query and re-renders on every data change. */
  const rail = useMediaQuery(RAIL);

  return (
    /* pl/pr-safe rather than a blanket inset: the top and bottom are handled
       by the surfaces that actually touch them, so a landscape notch does not
       steal height from a portrait screen. */
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
        {/* One nav slot, two presentations: the pull-down renders the same
            node so the phone can never drift from the desktop. */}
        {rail ? null : <MenuHandle>{nav}</MenuHandle>}

        <Outlet />
      </main>

      {palette}
    </div>
  );
}
