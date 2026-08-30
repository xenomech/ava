import { cn } from "@ava/ui";
import { Outlet } from "@tanstack/react-router";

import { useMediaQuery } from "@/shared/hooks/use-media-query";
import { AvaSocketProvider } from "@/shared/realtime";
import { AvaEvents } from "./ava-events";
import { CommandPalette } from "./command-palette";
import { MenuHandle } from "./menu-handle";
import { NavContent } from "./nav-content";

/** Where the rail replaces the pull-down. Matches the `md:` classes below. */
const RAIL = "(min-width: 768px)";

/**
 * No top bar. The window is the app: rooms down the side, the light filling
 * everything else. Anything configurable lives behind the account pill at the
 * foot of the sidebar rather than competing for space up top.
 *
 * On a phone the sidebar becomes a panel pulled down from the top edge, so the
 * only permanent chrome is a handle the width of two fingers.
 */
export function AppShell() {
  /* Mounted per form factor rather than hidden with CSS: a display:none rail
     still runs every query and re-renders on every device change. */
  const rail = useMediaQuery(RAIL);

  return (
    <AvaSocketProvider>
      <AvaEvents />

      {/* pl/pr-safe rather than a blanket inset: the top and bottom are handled
          by the surfaces that actually touch them, so a landscape notch does not
          steal height from a portrait screen. */}
      <div
        className={cn(
          "grid h-dvh bg-surface pl-safe pr-safe",
          "md:grid-cols-[208px_minmax(0,1fr)] xl:grid-cols-[262px_minmax(0,1fr)]",
        )}
      >
        {rail ? (
          <aside className="hidden min-h-0 bg-surface px-4 pb-8 pt-8 md:block xl:px-6 xl:pb-10 xl:pt-10">
            <NavContent />
          </aside>
        ) : null}

        <main
          className={cn(
            "relative m-2 min-h-0 min-w-0 overflow-y-auto rounded-2xl bg-bg",
            "mb-[max(0.5rem,env(safe-area-inset-bottom))] mt-[max(0.5rem,env(safe-area-inset-top))]",
          )}
        >
          {/* One nav list, two presentations: the pull-down renders the same
              component so the phone can never drift from the desktop. */}
          {rail ? null : (
            <MenuHandle>
              <NavContent />
            </MenuHandle>
          )}

          <Outlet />
        </main>

        <CommandPalette />
      </div>
    </AvaSocketProvider>
  );
}
