import { cn } from "@ava/ui";

// The dot is a light, lit in whatever colour the surrounding room is running at.
export function Wordmark({ className }: { className?: string }) {
  return (
    <div className={cn("flex items-center gap-2.5", className)}>
      <span
        aria-hidden
        className="size-2.5 shrink-0 rounded-full"
        style={{
          background: "var(--lit, #ffb463)",
          boxShadow: "0 0 16px 2px var(--lit, #ffb463)",
        }}
      />
      <span className="text-lead font-semibold lowercase tracking-snug">ava</span>
    </div>
  );
}
