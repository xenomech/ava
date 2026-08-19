import { Button } from "@ava/ui";
import { Link } from "@tanstack/react-router";

export function NotFound() {
  return (
    <main className="grid min-h-dvh place-items-center bg-bg p-6">
      <div className="grid w-full max-w-[400px] justify-items-start gap-3 rounded-xl border border-border bg-surface p-7">
        <span className="font-mono text-caption tracking-caps text-subtle uppercase">404</span>
        <h1 className="text-display font-semibold">Nothing here</h1>
        <p className="text-small text-muted">
          This page does not exist, or it belongs to a home you are no longer part of.
        </p>
        <Link to="/" className="mt-2 w-full">
          <Button variant="ghost" block>
            Back to the console
          </Button>
        </Link>
      </div>
    </main>
  );
}
