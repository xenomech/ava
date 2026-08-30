import { Button } from "@ava/ui";
import { Link } from "@tanstack/react-router";
import { RadioIcon } from "lucide-react";

/** Something was linked to that no longer exists. Deleted elsewhere, usually. */
export function Missing({ title, detail }: { title: string; detail: string }) {
  return (
    <div className="grid min-h-full place-items-center p-6">
      <div className="grid max-w-[360px] justify-items-center gap-3 text-center">
        <h1 className="text-title font-semibold">{title}</h1>
        <p className="text-small text-muted">{detail}</p>
        <Link to="/" className="mt-2">
          <Button>Back home</Button>
        </Link>
      </div>
    </div>
  );
}

/** Where `/` lands with no rooms, listing loose devices that would otherwise be unreachable. */
export function NoRooms({ devices }: { devices: { id: string; name: string }[] }) {
  return (
    <div className="grid min-h-full place-items-center p-6">
      <div className="grid w-full max-w-[420px] justify-items-center gap-3 text-center">
        <span className="grid size-12 place-items-center rounded-full bg-raised text-muted">
          <RadioIcon aria-hidden />
        </span>

        <h1 className="text-title font-semibold">No rooms yet</h1>

        <p className="text-small text-muted">
          Ava is organised by room — each one gets a page and a switch. Add your first from the
          menu.
        </p>

        {devices.length > 0 ? (
          <div className="mt-4 grid w-full gap-1.5 text-left">
            <span className="px-1 font-mono text-caption uppercase tracking-caps text-subtle">
              Not in a room
            </span>
            {devices.map((device) => (
              <Link
                key={device.id}
                to="/devices/$deviceId"
                params={{ deviceId: device.id }}
                className="truncate rounded-md border border-border bg-surface px-3 py-2.5 text-small font-medium hover:border-border-strong"
              >
                {device.name}
              </Link>
            ))}
          </div>
        ) : null}
      </div>
    </div>
  );
}

export function NoDevices({ hasHub }: { hasHub: boolean }) {
  return (
    <div className="grid min-h-full place-items-center p-6">
      <div className="grid max-w-[420px] justify-items-center gap-3 text-center">
        <span className="grid size-12 place-items-center rounded-full bg-raised text-muted">
          <RadioIcon aria-hidden />
        </span>

        <h1 className="text-title font-semibold">
          {hasHub ? "No devices yet" : "No hub connected"}
        </h1>

        <p className="text-small text-muted">
          {hasHub
            ? "Your hub is paired and sweeping the network every minute. Check your lights are switched on at the wall and joined to the same Wi-Fi — they will appear here on their own."
            : "Ava reaches your lights through a hub on your home network. Pair one and its devices show up here by themselves."}
        </p>

        {hasHub ? null : (
          <Link to="/activate" className="mt-2">
            <Button>Pair a hub</Button>
          </Link>
        )}
      </div>
    </div>
  );
}
