import { Button, Chip, Field, Input } from "@ava/ui";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { useSession } from "@/modules/auth";
import { Page, Section } from "@/shared/components/page";
import { updateCurrentTenant } from "../api";
import { tenantQueries } from "../queries";

export function SettingsPage() {
  const queryClient = useQueryClient();
  const session = useSession();
  const isAdmin = session.isAdmin;
  const tenant = useQuery(tenantQueries.current());

  const [name, setName] = useState("");

  useEffect(() => {
    if (tenant.data) setName(tenant.data.name);
  }, [tenant.data]);

  const save = useMutation({
    mutationFn: () => updateCurrentTenant(name),
    onSuccess: async () => {
      toast.success("Home updated");
      await queryClient.invalidateQueries({ queryKey: tenantQueries.all() });
    },
    onError: (error) => toast.error(isApiError(error) ? error.message : "Could not save"),
  });

  return (
    <Page
      description={tenant.data ? `ava.app/${tenant.data.slug}` : "Loading…"}
      actions={
        tenant.data ? (
          <Chip tone="muted" className="uppercase">
            {tenant.data.slug}
          </Chip>
        ) : null
      }
    >
      <Section title="Home" description="The name everyone in this home sees, on every device.">
        <form
          className="grid gap-4 p-5"
          onSubmit={(event) => {
            event.preventDefault();
            save.mutate();
          }}
        >
          <Field label="Name" error={isApiError(save.error) ? save.error.details?.name : undefined}>
            {(props) => (
              <Input
                {...props}
                required
                disabled={!isAdmin || tenant.isPending}
                value={name}
                onChange={(event) => setName(event.target.value)}
              />
            )}
          </Field>

          {isAdmin ? (
            <Button
              type="submit"
              className="justify-self-start"
              loading={save.isPending}
              disabled={!name || name === tenant.data?.name}
            >
              Save changes
            </Button>
          ) : (
            <p className="text-small text-muted">Only owners and admins can rename this home.</p>
          )}
        </form>
      </Section>

      <Section title="You">
        <dl className="grid gap-3 p-5 text-body">
          <Row label="Name" value={session.user?.name ?? "—"} />
          <Row label="Email" value={session.user?.email ?? "—"} />
          <Row
            label="Role in this home"
            value={<span className="capitalize">{session.tenant?.role ?? "—"}</span>}
          />
          <Row label="App version" value={<span className="font-mono">{__APP_VERSION__}</span>} />
        </dl>
      </Section>
    </Page>
  );
}

function Row({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex items-baseline justify-between gap-4">
      <dt className="text-small text-muted">{label}</dt>
      <dd className="min-w-0 truncate">{value}</dd>
    </div>
  );
}
