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

      <div className="grid h-dvh md:grid-cols-[208px_minmax(0,1fr)] xl:grid-cols-[262px_minmax(0,1fr)] bg-surface">
        <aside className="hidden min-h-0 bg-surface px-4 py-8 md:block xl:px-6 xl:py-10">
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
                "fixed left-3 top-3 z-sticky grid size-11 place-items-center rounded-full md:hidden",
                "border border-border bg-surface/80 text-fg backdrop-blur-sm",
              )}
            >
              <MenuIcon className="size-5" aria-hidden />
            </button>
          </DrawerTrigger>

          <DrawerContent
            grabber={false}
            className="inset-y-0 left-0 right-auto mt-0 h-full w-[280px] max-w-[85vw] rounded-none rounded-r-2xl border-r border-t-0"
          >
            <DrawerTitle className="sr-only">Menu</DrawerTitle>
            <DrawerDescription className="sr-only">Rooms and account</DrawerDescription>

            <div className="min-h-0 flex-1 px-4 py-8">
              <NavContent onNavigate={() => setMenuOpen(false)} />
            </div>
          </DrawerContent>
        </Drawer>

        <main className="min-h-0 min-w-0 overflow-y-auto bg-bg m-2 rounded-2xl">
          <Outlet />
        </main>

        <CommandPalette />
      </div>
    </AvaSocketProvider>
  );
}
