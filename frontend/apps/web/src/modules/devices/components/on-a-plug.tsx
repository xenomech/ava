import { cn } from "@ava/ui";
import { PlugZapIcon } from "lucide-react";

/** The badge for an appliance that is smart only by the plug it sits on. */
export function OnAPlug({ className }: { className?: string }) {
  return (
    <span
      title="On a smart plug"
      aria-label="On a smart plug"
      className={cn(
        "grid size-6 place-items-center rounded-full border border-border bg-surface text-subtle",
        className,
      )}
    >
      <PlugZapIcon className="size-3.5" aria-hidden />
    </span>
  );
}
