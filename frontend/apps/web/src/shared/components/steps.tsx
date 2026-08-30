import { cn } from "@ava/ui";

export function Steps({
  labels,
  at,
  className,
}: {
  labels: readonly string[];
  at: number;
  className?: string;
}) {
  return (
    <div className={cn("grid gap-2", className)}>
      <span className="font-mono text-caption text-subtle tabular">
        Step {at + 1} of {labels.length}
      </span>

      <div className="flex gap-1" aria-hidden>
        {labels.map((label, i) => (
          <span
            key={label}
            className={cn(
              "h-[3px] flex-1 rounded-[2px] transition-colors duration-300 ease-out",
              i <= at ? "bg-fg" : "bg-border",
            )}
          />
        ))}
      </div>
    </div>
  );
}
