import { Button, Chip, cn } from "@ava/ui";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { LogOutIcon } from "lucide-react";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { switchTenant, useSession, useSignOut } from "@/modules/auth";
import { tenantQueries } from "@/modules/tenant";
import { Page, Section } from "@/shared/components/page";

// Signing out and switching home both used to live in the top bar. With the bar
// gone they need a home of their own, and they are account concerns rather than
// home settings, so they get their own tab.
function Account() {
  const { user, tenant } = useSession();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const signOut = useSignOut();
  const homes = useQuery(tenantQueries.mine()).data ?? [];

  const switchTo = useMutation({
    mutationFn: switchTenant,
    onSuccess: () => {
      queryClient.clear();
      void navigate({ to: "/" });
    },
    onError: (error) => toast.error(isApiError(error) ? error.message : "Could not switch home"),
  });

  return (
    <Page>
      <Section title="You">
        <div className="grid gap-1 p-5">
          <p className="text-body font-semibold">{user?.name ?? "—"}</p>
          <p className="text-small text-muted">{user?.email}</p>
        </div>
      </Section>

      {homes.length > 1 ? (
        <Section title="Your homes" description="Switching signs you into that home.">
          <div className="grid gap-2 p-5">
            {homes.map((home) => (
              <button
                key={home.id}
                type="button"
                aria-current={home.slug === tenant?.slug}
                disabled={switchTo.isPending || home.slug === tenant?.slug}
                onClick={() => switchTo.mutate(home.slug)}
                className={cn(
                  "flex items-center justify-between gap-3 rounded-lg border border-border p-4 text-left",
                  "transition-colors duration-150 ease-out hover:border-border-strong",
                  "aria-[current=true]:border-fg",
                  "disabled:cursor-default",
                )}
              >
                <span className="min-w-0">
                  <span className="block truncate text-body font-semibold">{home.name}</span>
                  <span className="block text-caption text-subtle">{home.slug}</span>
                </span>
                <Chip tone="muted" className="uppercase">
                  {home.role}
                </Chip>
              </button>
            ))}
          </div>
        </Section>
      ) : null}

      <Section>
        <div className="flex items-center justify-between gap-4 p-5">
          <div className="grid gap-1">
            <p className="text-small font-semibold">Sign out</p>
            <p className="text-caption text-muted">Ends this session on this device.</p>
          </div>
          <Button
            variant="ghost"
            loading={signOut.isPending}
            onClick={() => signOut.mutate()}
            className="shrink-0 gap-2"
          >
            <LogOutIcon aria-hidden />
            Sign out
          </Button>
        </div>
      </Section>
    </Page>
  );
}

export const Route = createFileRoute("/_protected/settings/account")({
  component: Account,
});
