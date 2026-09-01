import type { ApiTokenDto } from "@ava/contracts";
import { Button, Chip, cn } from "@ava/ui";

import { scopeLabel } from "../lib/scopes";

function when(iso: string | null): string {
  if (!iso) return "never";

  return new Date(iso).toLocaleDateString(undefined, {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
}

// A token's state is one of three things, so it is said once rather than inferred from three fields.
function status(
  token: ApiTokenDto,
): { label: string; tone: "muted" | "warning" | "danger" } | null {
  if (token.revoked_at) return { label: "Revoked", tone: "danger" };
  if (token.expires_at && new Date(token.expires_at) < new Date()) {
    return { label: "Expired", tone: "warning" };
  }

  return null;
}

export function TokenRow({
  token,
  busy,
  onRevoke,
  onDelete,
}: {
  token: ApiTokenDto;
  busy: boolean;
  onRevoke: () => void;
  onDelete: () => void;
}) {
  const state = status(token);
  const dead = state !== null;

  return (
    <li className={cn("grid gap-3 p-5", dead && "opacity-60")}>
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0">
          <p className="flex items-center gap-2 truncate text-body font-medium">
            {token.name}
            {state ? <Chip tone={state.tone}>{state.label}</Chip> : null}
          </p>
          <p className="mt-1 text-caption text-muted">
            Last used {when(token.last_used_at)} · Expires {when(token.expires_at)}
          </p>
        </div>

        <div className="flex shrink-0 items-center gap-2">
          {dead ? null : (
            <Button variant="ghost" size="sm" disabled={busy} onClick={onRevoke}>
              Revoke
            </Button>
          )}
          <Button variant="ghost" size="sm" disabled={busy} onClick={onDelete}>
            Delete
          </Button>
        </div>
      </div>

      <div className="flex flex-wrap gap-1.5">
        {token.scopes.map((scope) => (
          <Chip key={scope} tone="muted">
            {scopeLabel(scope)}
          </Chip>
        ))}
      </div>
    </li>
  );
}
