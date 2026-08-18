import { env } from "@ava/env/web";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@ava/ui/components/card";
import { Skeleton } from "@ava/ui/components/skeleton";
import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/")({
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
      <h1 className="font-semibold text-2xl tracking-tight">Ava</h1>
      <p className="mt-1 text-muted-foreground text-sm">
        Multi-tenant control plane for hubs and the devices behind them.
      </p>

      <Card className="mt-8 max-w-md">
        <CardHeader>
          <CardTitle>API</CardTitle>
          <CardDescription>{env.VITE_API_URL}</CardDescription>
        </CardHeader>
        <CardContent>
          {health.state === "loading" && <Skeleton className="h-5 w-24" />}
          {health.state === "up" && (
            <span className="inline-flex items-center gap-2 text-sm">
              <span className="size-2 rounded-full bg-emerald-500" />
              {health.status}
            </span>
          )}
          {health.state === "down" && (
            <span className="inline-flex items-center gap-2 text-sm text-muted-foreground">
              <span className="size-2 rounded-full bg-destructive" />
              unreachable
            </span>
          )}
        </CardContent>
      </Card>
    </main>
  );
}
