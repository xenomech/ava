import { cn } from "@ava/ui";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { CheckIcon } from "lucide-react";
import type { CSSProperties } from "react";

import { isApiError } from "@/config/http/request";
import { activateHub, hubQueries } from "@/modules/hub";
import { CodeField } from "./code-field";

// This step pairs on its own; Continue only asks the server to confirm a hub arrived.
export function HubStep({
  code,
  onCodeChange,
  error,
}: {
  code: string;
  onCodeChange: (value: string) => void;
  error?: string;
}) {
  const queryClient = useQueryClient();
  const hubs = useQuery(hubQueries.list());
  const paired = hubs.data ?? [];

  const pair = useMutation({
    mutationFn: (value: string) => activateHub({ user_code: value.trim().toUpperCase() }),
    onSuccess: async () => {
      onCodeChange("");
      await queryClient.invalidateQueries({ queryKey: hubQueries.all() });
    },
  });

  const pairError = isApiError(pair.error)
    ? pair.error.message
    : pair.error
      ? "Could not pair that hub"
      : undefined;

  if (paired.length > 0) {
    return (
      <div className="grid gap-5">
        {paired.map((hub) => (
          <div key={hub.id} className="flex items-center gap-4">
            <Found />
            <span className="min-w-0">
              <b className="block truncate text-title font-semibold">{hub.name}</b>
              <span className="block text-small text-muted">
                Connected · looking for your lights now
              </span>
            </span>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className="grid gap-8 sm:grid-cols-[auto_minmax(0,1fr)] sm:items-center sm:gap-10">
      <Searching pairing={pair.isPending} />

      <div className="grid gap-4">
        <CodeField
          label="Pairing code"
          value={code}
          disabled={pair.isPending}
          error={error ?? pairError}
          onChange={onCodeChange}
          // A complete code pairs itself; hunting for a button is a step nobody needs.
          onComplete={(value) => pair.mutate(value)}
        />

        <p className="text-small text-subtle">Your hub shows this when it starts up.</p>
      </div>
    </div>
  );
}

// Three offset rings, so one is always mid-flight and it reads as listening.
const RINGS = [0, 0.93, 1.86];

function Searching({ pairing }: { pairing: boolean }) {
  return (
    <div
      aria-hidden
      className="relative mx-auto grid size-[132px] shrink-0 place-items-center sm:mx-0"
    >
      {RINGS.map((delay) => (
        <span
          key={delay}
          className="absolute size-full rounded-full border border-[var(--lit)]"
          style={{
            animation: `sweep 2.8s var(--ease-out-quart) ${delay}s infinite`,
            // Tighten the pulse while a code is actually being checked.
            animationDuration: pairing ? "1.1s" : undefined,
          }}
        />
      ))}

      <span
        className="size-3 rounded-full bg-[var(--lit)]"
        style={{ boxShadow: "0 0 20px 4px var(--lit)" }}
      />
    </div>
  );
}

function Found() {
  return (
    <span
      aria-hidden
      className={cn(
        "grid size-12 shrink-0 place-items-center rounded-full",
        "animate-scale-in border border-[var(--lit)] text-[var(--lit)]",
      )}
      style={{ boxShadow: "0 0 28px -4px var(--lit)" } as CSSProperties}
    >
      <CheckIcon className="size-6" />
    </span>
  );
}
