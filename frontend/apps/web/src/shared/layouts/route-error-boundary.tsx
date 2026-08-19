import { Button } from "@ava/ui";
import { useRouter, type ErrorComponentProps } from "@tanstack/react-router";

import { ApiError } from "@/config/http/types";

export function RouteErrorBoundary({ error, reset }: ErrorComponentProps) {
  const router = useRouter();

  const message =
    error instanceof ApiError ? error.message : "Something went wrong loading this page.";

  return (
    <main className="grid min-h-dvh place-items-center bg-bg p-6">
      <div className="grid w-full max-w-[440px] gap-3 rounded-xl border border-border bg-surface p-7">
        <h1 className="text-title font-semibold">That did not load</h1>
        <p className="text-small text-muted">{message}</p>

        <div className="mt-2 flex gap-2">
          <Button
            onClick={() => {
              reset();
              void router.invalidate();
            }}
          >
            Try again
          </Button>
          <Button variant="ghost" onClick={() => window.location.reload()}>
            Reload
          </Button>
        </div>

        {import.meta.env.DEV && error instanceof Error ? (
          <pre className="mt-2 max-w-full overflow-x-auto rounded-sm bg-raised p-4 text-left font-mono text-caption text-muted">
            {error.stack ?? error.message}
          </pre>
        ) : null}
      </div>
    </main>
  );
}
