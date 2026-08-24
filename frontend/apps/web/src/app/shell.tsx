import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerTitle,
  DrawerTrigger,
  cn,
} from "@ava/ui";
import { Outlet } from "@tanstack/react-router";
import { MenuIcon } from "lucide-react";
import { useState } from "react";

import { AvaSocketProvider } from "@/shared/realtime";
import { AvaEvents } from "./ava-events";
import { CommandPalette } from "./command-palette";
import { NavContent } from "./nav-content";

/**
 * No top bar. The window is the app: rooms down the side, the light filling
 * everything else. Anything configurable lives behind the account pill at the
 * foot of the sidebar rather than competing for space up top.
 */
export function AppShell() {
  const [menuOpen, setMenuOpen] = useState(false);

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
        <aside className="hidden min-h-0 bg-surface px-4 pb-8 pt-8 md:block xl:px-6 xl:pb-10 xl:pt-10">
          <NavContent />
        </aside>

        {/* One nav list, two presentations: the drawer renders the same
            component so the phone can never drift from the desktop. */}
        <Drawer open={menuOpen} onOpenChange={setMenuOpen} direction="left">
          <DrawerTrigger asChild>
            <button
              type="button"
              aria-label="Open menu"
              className={cn(
                /* 44px is the smallest target a thumb finds reliably, and the
                   inset keeps it clear of the notch on a standalone install. */
                "fixed left-3 z-sticky grid size-11 place-items-center rounded-full md:hidden",
                "top-[max(0.75rem,calc(env(safe-area-inset-top)+0.25rem))]",
                "border border-border bg-surface/80 text-fg backdrop-blur-sm",
              )}
            >
              <MenuIcon className="size-5" aria-hidden />
            </button>
          </DrawerTrigger>

          <DrawerContent
            grabber={false}
            className="inset-y-0 left-0 right-auto mt-0 h-full w-70 max-w-[85vw] rounded-none rounded-r-2xl border-r border-t-0"
          >
            <DrawerTitle className="sr-only">Menu</DrawerTitle>
            <DrawerDescription className="sr-only">Rooms and account</DrawerDescription>

            <div className="min-h-0 flex-1 overflow-y-auto px-4 pb-[max(2rem,env(safe-area-inset-bottom))] pt-[max(2rem,env(safe-area-inset-top))]">
              <NavContent onNavigate={() => setMenuOpen(false)} />
            </div>
          </DrawerContent>
        </Drawer>

        <main
          className={cn(
            "m-2 min-h-0 min-w-0 overflow-y-auto rounded-2xl bg-bg",
            "mb-[max(0.5rem,env(safe-area-inset-bottom))] mt-[max(0.5rem,env(safe-area-inset-top))]",
          )}
        >
          <Outlet />
        </main>

        <CommandPalette />
      </div>
    </AvaSocketProvider>
  );
}
