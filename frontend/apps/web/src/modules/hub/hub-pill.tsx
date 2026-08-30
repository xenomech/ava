import { cn } from "@ava/ui";
import type { HubDto } from "@ava/contracts";

/**
 * How long a hub can go quiet before the room stops believing in it.
 *
 * Matches the grace the API applies when it decides whether to accept a
 * command, so the pill and the thing it is describing agree. A hub reports
 * every thirty seconds, which gives it several missed beats before anyone
 * panics.
 */
const GRACE_MS = 3 * 60 * 1000;

export type HubHealth = "connected" | "unreachable" | "offline";

/**
 * What the room can actually expect of its hub.
 *
 * Three states rather than two, because the middle one is real and cost a whole
 * afternoon to find: a hub can be perfectly alive — heartbeating, syncing
 * devices, reporting for duty over HTTP — while the channel that carries
 * commands is dead. Everything looked healthy and nothing worked. Flattening
 * that into "offline" would be a lie, and flattening it into "connected" is the
 * lie we already lived through.
 */
export function hubHealth(hub: HubDto, now = Date.now()): HubHealth {
  if (hub.status === "revoked") return "offline";
  if (hub.online) return "connected";

  const seen = hub.last_seen_at ? Date.parse(hub.last_seen_at) : Number.NaN;

  /* Heard from, but not reachable: it is there and something between here and
     it is broken. That is worth telling apart from silence. */
  return Number.isNaN(seen) || now - seen > GRACE_MS ? "offline" : "unreachable";
}

const TONE: Record<HubHealth, { dot: string; says: string }> = {
  connected: { dot: "bg-success", says: "connected" },
  unreachable: { dot: "bg-warning", says: "cannot be reached" },
  offline: { dot: "bg-danger", says: "offline" },
};

/**
 * The hub, named, with a light on it.
 *
 * Sits beside the room's name because it qualifies the room: everything on this
 * screen is a promise about a set of bulbs, and this says whether the promise
 * can currently be kept.
 */
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
          /* Only the state that wants attention gets any. A steady green ring
             on a screen you look at every day is just decoration. */
          health === "unreachable" && "animate-pulse",
        )}
      />
      <span className="min-w-0 truncate">{hub.name}</span>
      <span className="sr-only">— {tone.says}</span>
    </span>
  );
}
