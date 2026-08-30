import { cn } from "@ava/ui";
import { Link } from "@tanstack/react-router";

import { useSession } from "../hooks/use-session";

// One target for the account and everything configurable behind it.
export function UserPill() {
  const { user, tenant } = useSession();
  const initial = (user?.name ?? user?.email ?? "?").trim().slice(0, 1).toUpperCase();

  return (
    <Link
      to="/settings"
      className={cn(
        "mt-4 flex min-h-[56px] items-center gap-3 rounded-[14px] p-2.5",
        "text-muted transition-colors duration-[180ms] ease-out-soft",
        "hover:bg-raised hover:text-fg",
        "aria-[current=page]:bg-raised aria-[current=page]:text-fg",
      )}
    >
      <span
        aria-hidden
        className="grid size-9 shrink-0 place-items-center rounded-full bg-raised text-small font-semibold text-fg"
      >
        {initial}
      </span>

      <span className="min-w-0 flex-1">
        <span className="block truncate text-[0.8125rem] font-medium text-fg">
          {user?.name ?? "Account"}
        </span>
        <span className="block truncate font-mono text-micro text-subtle">
          {tenant?.name ?? "ava"}
        </span>
      </span>
    </Link>
  );
}
