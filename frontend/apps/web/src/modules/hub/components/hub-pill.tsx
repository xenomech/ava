import { cn } from "@ava/ui";
import type { HubDto } from "@ava/contracts";

// How long a hub can go quiet; matches the grace the API applies to commands.
const GRACE_MS = 3 * 60 * 1000;

export type HubHealth = "connected" | "unreachable" | "offline";

/** What the room can expect of its hub; a hub can be alive but unable to command. */
export function hubHealth(hub: HubDto, now = Date.now()): HubHealth {
  if (hub.status === "revoked") return "offline";
  if (hub.online) return "connected";

  const seen = hub.last_seen_at ? Date.parse(hub.last_seen_at) : Number.NaN;

  // Heard from but not reachable is worth telling apart from silence.
  return Number.isNaN(seen) || now - seen > GRACE_MS ? "offline" : "unreachable";
}

const TONE: Record<HubHealth, { dot: string; says: string }> = {
  connected: { dot: "bg-success", says: "connected" },
  unreachable: { dot: "bg-warning", says: "cannot be reached" },
  offline: { dot: "bg-danger", says: "offline" },
};

/** The hub, named, with a light on it; it qualifies the room beside it. */
export function HubPill({ hub, className }: { hub: HubDto; className?: string }) {
  const health = hubHealth(hub);
  const tone = TONE[health];

  return (
    <span
      title={`${hub.name} — ${tone.says}`}
      className={cn(
        "inline-flex min-w-0 shrink-0 items-center gap-1.5 rounded-full border border-border",
        "bg-surface py-1 pl-2 pr-2.5 text-caption text-muted",
        className,
      )}
    >
      <span
        aria-hidden
        className={cn(
          "size-1.5 shrink-0 rounded-full transition-colors duration-300 ease-out",
          tone.dot,
          // Only the state that wants attention gets any; a steady ring is decoration.
          health === "unreachable" && "animate-pulse",
        )}
      />
      <span className="min-w-0 truncate">{hub.name}</span>
      <span className="sr-only">— {tone.says}</span>
    </span>
  );
}
