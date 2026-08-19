import { env } from "@ava/env/web";
import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/_protected/")({
  component: HomeComponent,
});

type Health = { state: "loading" } | { state: "up"; status: string } | { state: "down"; reason: string };

function useApiHealth() {
  const [health, setHealth] = useState<Health>({ state: "loading" });

  useEffect(() => {
    const controller = new AbortController();

    fetch(`${env.VITE_API_URL}/health`, { signal: controller.signal })
      .then(async (response) => {
        const body = await response.json();
        setHealth({ state: "up", status: body?.data?.status ?? "ok" });
      })
      .catch((error: unknown) => {
        if (controller.signal.aborted) {
          return;
        }

        setHealth({ state: "down", reason: error instanceof Error ? error.message : "unreachable" });
      });

    return () => controller.abort();
  }, []);

  return health;
}

function HomeComponent() {
  const health = useApiHealth();

  return (
    <main className="mx-auto w-full max-w-5xl px-4 py-10">
      <h1 className="text-display font-semibold">Ava</h1>
      <p className="mt-1 text-small text-muted">
        Multi-tenant control plane for hubs and the devices behind them.
      </p>

      <section className="mt-8 grid max-w-md gap-3 rounded-lg border border-border bg-surface p-5">
        <div className="grid gap-1">
          <h2 className="text-lead font-semibold">API</h2>
          <p className="font-mono text-caption text-subtle">{env.VITE_API_URL}</p>
        </div>

        {health.state === "loading" ? (
          <div className="h-5 w-24 animate-pulse rounded-xs bg-raised" />
        ) : null}

        {health.state === "up" ? (
          <span className="inline-flex items-center gap-2 text-small">
            <span className="size-2 rounded-full bg-success" />
            {health.status}
          </span>
        ) : null}

        {health.state === "down" ? (
          <span className="inline-flex items-center gap-2 text-small text-muted">
            <span className="size-2 rounded-full bg-danger" />
            unreachable
          </span>
        ) : null}
      </section>
    </main>
  );
}
