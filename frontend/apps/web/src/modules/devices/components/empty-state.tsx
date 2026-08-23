import { Button } from "@ava/ui";
import { Link } from "@tanstack/react-router";
import { RadioIcon } from "lucide-react";

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
