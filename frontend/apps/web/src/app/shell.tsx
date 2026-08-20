import {
  Button,
  Chip,
  cn,
  Drawer,
  DrawerBody,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@ava/ui";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link, Outlet, useNavigate } from "@tanstack/react-router";
import {
  ChevronDownIcon,
  HomeIcon,
  LayoutGridIcon,
  LogOutIcon,
  MoonIcon,
  PlusIcon,
  RadioIcon,
  SettingsIcon,
  SunIcon,
  UsersIcon,
} from "lucide-react";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { switchTenant, useSession, useSignOut } from "@/modules/auth";
import { tenantQueries } from "@/modules/tenant";
import { useTheme } from "@/shared/components/theme-provider";
import { CommandPalette } from "./command-palette";
import { useAvaStream } from "./use-ava-stream";

const NAV = [
  { to: "/", label: "Console", icon: HomeIcon },
  { to: "/rooms", label: "Rooms", icon: LayoutGridIcon },
  { to: "/activate", label: "Hubs", icon: RadioIcon },
  { to: "/settings/members", label: "People", icon: UsersIcon },
  { to: "/settings", label: "Settings", icon: SettingsIcon },
] as const;

const MOBILE_NAV = NAV.filter((item) => item.to !== "/settings/members");

export function AppShell() {
  const { user, tenant } = useSession();

  useAvaStream();

  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const homes = useQuery(tenantQueries.mine()).data ?? [];

  const { resolvedTheme, setTheme } = useTheme();
  const signOut = useSignOut();

  const switchTo = useMutation({
    mutationFn: switchTenant,
    onSuccess: () => {
      queryClient.clear();
      void navigate({ to: "/" });
    },
    onError: (error) => toast.error(isApiError(error) ? error.message : "Could not switch home"),
  });

  return (
    <div className="grid h-dvh grid-rows-[auto_minmax(0,1fr)_auto] bg-bg">
      <header className="flex items-center gap-3 border-b border-border px-4 py-2.5">
        <Drawer>
          <DrawerTrigger asChild>
            <Button variant="quiet" size="sm" className="gap-2 px-2 text-fg">
              <span className="grid size-6 place-items-center rounded-xs bg-accent text-caption font-bold text-accent-fg">
                a
              </span>
              {tenant?.name ?? "ava"}
              <ChevronDownIcon />
            </Button>
          </DrawerTrigger>

          <DrawerContent>
            <DrawerHeader>
              <div>
                <DrawerTitle>Your homes</DrawerTitle>
                <DrawerDescription>Switching signs you into that home</DrawerDescription>
              </div>
            </DrawerHeader>

            <DrawerBody>
              <div className="grid gap-2 pb-4">
                {homes.map((home) => (
                  <DrawerClose asChild key={home.id}>
                    <button
                      type="button"
                      aria-current={home.slug === tenant?.slug}
                      disabled={switchTo.isPending}
                      onClick={() => switchTo.mutate(home.slug)}
                      className={cn(
                        "flex items-center justify-between gap-3 rounded-lg border border-border p-4 text-left",
                        "transition-colors duration-150 ease-out hover:border-border-strong",
                        "aria-[current=true]:border-fg",
                      )}
                    >
                      <span>
                        <span className="block text-lead font-semibold">{home.name}</span>
                        <span className="block text-caption text-subtle">{home.slug}</span>
                      </span>
                      <Chip tone="muted" className="uppercase">
                        {home.role}
                      </Chip>
                    </button>
                  </DrawerClose>
                ))}

                {homes.length === 0 ? (
                  <p className="py-6 text-center text-small text-muted">No other homes yet.</p>
                ) : null}
              </div>
            </DrawerBody>
          </DrawerContent>
        </Drawer>

        <div className="flex-1" />

        <kbd className="hidden items-center gap-1 rounded-xs border border-border px-1.5 py-0.5 font-mono text-caption text-subtle md:inline-flex">
          ⌘K
        </kbd>

        <span className="hidden text-small text-muted sm:inline">{user?.email}</span>

        <Button
          variant="ghost"
          size="icon-sm"
          aria-label={resolvedTheme === "dark" ? "Switch to light" : "Switch to dark"}
          onClick={() => setTheme(resolvedTheme === "dark" ? "light" : "dark")}
        >
          {resolvedTheme === "dark" ? <SunIcon aria-hidden /> : <MoonIcon aria-hidden />}
        </Button>

        <Button
          variant="ghost"
          size="icon-sm"
          aria-label="Sign out"
          loading={signOut.isPending}
          onClick={() => signOut.mutate()}
        >
          <LogOutIcon aria-hidden />
        </Button>
      </header>

      <div className="grid min-h-0 md:grid-cols-[64px_minmax(0,1fr)] lg:grid-cols-[208px_minmax(0,1fr)]">
        <nav
          aria-label="Sections"
          className="hidden min-h-0 flex-col gap-1 overflow-y-auto border-r border-border p-3 md:flex"
        >
          {NAV.map(({ to, label, icon: Icon }) => (
            <Link
              key={to}
              to={to}
              activeOptions={{ exact: true }}
              className={cn(
                "flex h-10 items-center justify-center gap-3 rounded-sm px-3 text-small font-medium lg:justify-start",
                "transition-colors duration-150 ease-out",
                "text-muted hover:bg-surface hover:text-fg",
                "aria-[current=page]:bg-raised aria-[current=page]:text-fg",
              )}
            >
              <Icon className="size-[18px] shrink-0" aria-hidden />
              <span className="hidden lg:inline">{label}</span>
            </Link>
          ))}

          <Link
            to="/activate"
            className={cn(
              "mt-2 flex h-10 items-center justify-center gap-3 rounded-sm px-3 text-small font-medium lg:justify-start",
              "border border-border-strong text-fg transition-colors duration-150 ease-out hover:bg-surface",
            )}
          >
            <PlusIcon className="size-[18px] shrink-0" aria-hidden />
            <span className="hidden lg:inline">Add hub</span>
          </Link>
        </nav>

        <main className="min-w-0 overflow-y-auto">
          <Outlet />
        </main>
      </div>

      <nav
        aria-label="Main"
        className="grid grid-cols-5 items-stretch border-t border-border pb-[env(safe-area-inset-bottom)] md:hidden"
      >
        {MOBILE_NAV.map(({ to, label, icon: Icon }) => (
          <Link
            key={to}
            to={to}
            activeOptions={{ exact: true }}
            className={cn(
              "grid h-14 place-items-center gap-1 text-caption",
              "transition-colors duration-150 ease-out",
              "text-subtle aria-[current=page]:text-fg",
            )}
          >
            <Icon className="size-[19px]" aria-hidden />
            <span className="leading-none">{label}</span>
          </Link>
        ))}

        <Link
          to="/activate"
          className="grid h-14 place-items-center gap-1 text-caption text-subtle"
        >
          <span className="grid size-[19px] place-items-center rounded-xs bg-accent text-accent-fg">
            <PlusIcon className="size-3.5" aria-hidden />
          </span>
          <span className="leading-none">Add</span>
        </Link>
      </nav>

      <CommandPalette />
    </div>
  );
}
