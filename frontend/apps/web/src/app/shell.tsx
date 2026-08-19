import { Button, Chip, buttonVariants } from "@ava/ui";
import { Link, Outlet } from "@tanstack/react-router";

import { useSession, useSignOut } from "@/modules/auth";
import { ModeToggle } from "@/shared/components/mode-toggle";

export function AppShell() {
  return (
    <div className="grid h-svh grid-rows-[auto_1fr] bg-bg">
      <Header />
      <div className="overflow-y-auto">
        <Outlet />
      </div>
    </div>
  );
}

function Header() {
  const session = useSession();
  const signOut = useSignOut();

  return (
    <header className="border-b border-border">
      <div className="mx-auto flex w-full max-w-5xl items-center justify-between px-4 py-3">
        <div className="flex items-center gap-6">
          <Link to="/" className="text-lead font-semibold">
            Ava
          </Link>
          {session.tenant ? (
            <span className="flex items-center gap-2.5 text-small text-muted">
              {session.tenant.name}
              <Chip tone="muted" className="px-2 py-0.5 text-caption capitalize">
                {session.tenant.role}
              </Chip>
            </span>
          ) : null}
        </div>

        <div className="flex items-center gap-3">
          {session.isAuthenticated ? (
            <>
              <span className="text-small text-muted">{session.user?.email}</span>
              <Button
                variant="ghost"
                size="sm"
                disabled={signOut.isPending}
                onClick={() => signOut.mutate()}
              >
                Sign out
              </Button>
            </>
          ) : session.isLoading ? null : (
            <Link to="/auth/login" className={buttonVariants({ variant: "ghost", size: "sm" })}>
              Sign in
            </Link>
          )}
          <ModeToggle />
        </div>
      </div>
    </header>
  );
}
