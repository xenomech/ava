import type { ApiTokenDto, CreateApiTokenRequest } from "@ava/contracts";
import { Confirm } from "@ava/ui";
import { useState } from "react";

import { Page, Section } from "@/shared/components/page";
import { TokenForm } from "../components/token-form";
import { TokenReveal } from "../components/token-reveal";
import { TokenRow } from "../components/token-row";
import { useTokenActions, useTokens } from "../hooks/use-tokens";

export function TokensPage() {
  const { tokens, isPending } = useTokens();
  const { create, revoke, remove } = useTokenActions();

  // Held in state, not refetched: the server will never hand this value over again.
  const [issued, setIssued] = useState<string | null>(null);
  const [doomed, setDoomed] = useState<ApiTokenDto | null>(null);

  const onCreate = (body: CreateApiTokenRequest) =>
    create.mutate(body, { onSuccess: (created) => setIssued(created.value) });

  return (
    <Page description="Tokens let a script or a shortcut act for you, without your password.">
      {issued ? (
        <Section>
          <div className="p-5">
            <TokenReveal value={issued} onDone={() => setIssued(null)} />
          </div>
        </Section>
      ) : (
        <Section title="New token" description="Shown once, so keep it somewhere safe.">
          <TokenForm busy={create.isPending} onCreate={onCreate} />
        </Section>
      )}

      <Section title="Your tokens">
        {isPending ? (
          <p className="p-5 text-small text-muted">Loading tokens…</p>
        ) : tokens.length === 0 ? (
          <p className="p-5 text-small text-muted">
            No tokens yet. Create one above to drive Ava from a shortcut.
          </p>
        ) : (
          <ul className="divide-y divide-border">
            {tokens.map((token) => (
              <TokenRow
                key={token.id}
                token={token}
                busy={revoke.isPending || remove.isPending}
                onRevoke={() => revoke.mutate(token.id)}
                onDelete={() => setDoomed(token)}
              />
            ))}
          </ul>
        )}
      </Section>

      <Confirm
        open={doomed !== null}
        onOpenChange={(open) => !open && setDoomed(null)}
        title={`Delete ${doomed?.name ?? "this token"}?`}
        description="Anything using it stops working immediately, and the record goes with it."
        confirmLabel="Delete token"
        tone="danger"
        onConfirm={() => {
          if (doomed) remove.mutate(doomed.id);
          setDoomed(null);
        }}
      />
    </Page>
  );
}
