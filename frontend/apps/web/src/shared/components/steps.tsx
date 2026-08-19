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
    <div className={cn("grid gap-2.5", className)}>
      <div className="flex items-baseline justify-between gap-3">
        <span className="text-small font-semibold">{labels[at]}</span>
        <span className="font-mono text-caption text-subtle tabular">
          {at + 1} of {labels.length}
        </span>
      </div>

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
