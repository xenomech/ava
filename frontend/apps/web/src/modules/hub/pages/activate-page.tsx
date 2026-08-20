import { Button, Chip, Field, Input } from "@ava/ui";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useSearch } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { useSession } from "@/modules/auth";
import { Page, Section } from "@/shared/components/page";
import { activateHub, removeHub } from "../api";
import { hubQueries } from "../queries";

export function ActivatePage() {
  const queryClient = useQueryClient();
  const { isAdmin } = useSession();
  const search = useSearch({ strict: false }) as { code?: string };

  const hubs = useQuery(hubQueries.list());
  const [code, setCode] = useState("");

  useEffect(() => {
    if (search.code) setCode(search.code);
  }, [search.code]);

  const activate = useMutation({
    mutationFn: () => activateHub({ user_code: code.trim().toUpperCase() }),
    onSuccess: async (hub) => {
      toast.success(`${hub.name} is paired`);
      setCode("");
      await queryClient.invalidateQueries({ queryKey: hubQueries.all() });
    },
    onError: (error) => toast.error(isApiError(error) ? error.message : "Could not pair that hub"),
  });

  const revoke = useMutation({
    mutationFn: removeHub,
    onSuccess: async () => {
      toast.success("Hub removed");
      await queryClient.invalidateQueries({ queryKey: hubQueries.all() });
    },
    onError: (error) => toast.error(isApiError(error) ? error.message : "Could not remove the hub"),
  });

  return (
    <Page title="Hubs" description="Pair the box that talks to your devices.">
      <Section title="Add a hub" description="Start Ava on the hub and enter the code it prints.">
        <form
          className="grid gap-4 p-5 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-end"
          onSubmit={(event) => {
            event.preventDefault();
            activate.mutate();
          }}
        >
          <Field
            label="Pairing code"
            hint="Nine characters, like ABCD-1234."
            error={isApiError(activate.error) ? activate.error.details?.user_code : undefined}
          >
            {(props) => (
              <Input
                {...props}
                required
                autoComplete="off"
                spellCheck={false}
                placeholder="ABCD-1234"
                className="font-mono tracking-caps uppercase"
                value={code}
                onChange={(event) => setCode(event.target.value)}
              />
            )}
          </Field>

          <Button type="submit" loading={activate.isPending} disabled={!code.trim()}>
            Pair hub
          </Button>
        </form>
      </Section>

      <Section>
        {hubs.isPending ? (
          <p className="p-5 text-small text-muted">Loading hubs…</p>
        ) : hubs.isError ? (
          <p role="alert" className="p-5 text-small text-danger">
            {isApiError(hubs.error) ? hubs.error.message : "Could not load hubs"}
          </p>
        ) : (
          <ul className="divide-y divide-border">
            {hubs.data.map((hub) => (
              <li key={hub.id} className="flex items-center justify-between gap-4 p-5">
                <div className="min-w-0">
                  <p className="truncate text-body font-medium">{hub.name}</p>
                  <p className="truncate font-mono text-caption text-subtle">
                    {hub.last_seen_at
                      ? `last seen ${new Date(hub.last_seen_at).toLocaleString()}`
                      : "never seen"}
                  </p>
                </div>

                <div className="flex shrink-0 items-center gap-2">
                  <Chip tone={hub.status === "active" ? "success" : "muted"} className="capitalize">
                    {hub.status}
                  </Chip>

                  {isAdmin ? (
                    <Button
                      variant="ghost"
                      size="sm"
                      loading={revoke.isPending && revoke.variables === hub.id}
                      onClick={() => revoke.mutate(hub.id)}
                    >
                      Remove
                    </Button>
                  ) : null}
                </div>
              </li>
            ))}

            {hubs.data.length === 0 ? (
              <li className="p-8 text-center text-small text-muted">
                No hubs yet. Run Ava on your Raspberry Pi and its code will appear.
              </li>
            ) : null}
          </ul>
        )}
      </Section>
    </Page>
  );
}
