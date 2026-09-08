import { cn } from "@ava/ui";
import { Link, Outlet, createFileRoute } from "@tanstack/react-router";

// Everything configurable lives behind one door, as tabs, rather than competing with the rooms.
const TABS = [
  { to: "/settings", label: "General", exact: true },
  { to: "/settings/people", label: "People", exact: false },
  { to: "/settings/hubs", label: "Hubs", exact: false },
  { to: "/settings/appearance", label: "Appearance", exact: false },
  { to: "/settings/account", label: "Account", exact: false },
  { to: "/settings/tokens", label: "Tokens", exact: false },
] as const;

function SettingsLayout() {
  return (
    // Heading and tabs bring their own container, so they do not stack padding on the tab's page.
    <div className="pt-3 md:pt-0">
      <div className="mx-auto grid w-full max-w-[720px] gap-4 px-5 pt-6 sm:px-8 sm:pt-8">
        <h1 className="text-display font-semibold">Settings</h1>

        {/* Bled to the edge, so an overflowing tab is cut by the screen, not by inner padding. */}
        <nav
          aria-label="Settings"
          className={cn(
            "no-scrollbar -mx-5 flex snap-x gap-1.5 overflow-x-auto px-5 sm:-mx-8 sm:px-8",
            "scroll-px-5 sm:scroll-px-8",
          )}
        >
          {TABS.map((tab) => (
            <Link
              key={tab.to}
              to={tab.to}
              activeOptions={{ exact: tab.exact }}
              className={cn(
                "flex min-h-11 shrink-0 snap-start items-center rounded-full border border-border px-4",
                "text-small font-medium text-muted",
                "transition-colors duration-150 ease-out hover:text-fg",
                "aria-[current=page]:border-fg aria-[current=page]:text-fg",
                "[@media(hover:hover)]:min-h-9",
              )}
            >
              {tab.label}
            </Link>
          ))}
        </nav>
      </div>

      <Outlet />
    </div>
  );
}

export const Route = createFileRoute("/_protected/settings")({
  component: SettingsLayout,
});
