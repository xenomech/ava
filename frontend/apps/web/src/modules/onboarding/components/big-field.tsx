import { cn } from "@ava/ui";
import { useId, type ComponentProps } from "react";

// A boxed input would read as a form. At this size the text is the interface,
// so the only chrome is a rule underneath that lights up when focused — in the
// room's own colour, like everything else here.
export function BigField({
  label,
  error,
  hint,
  className,
  ...props
}: ComponentProps<"input"> & { label: string; error?: string; hint?: string }) {
  const id = useId();
  const describedBy = error ? `${id}-error` : hint ? `${id}-hint` : undefined;

  return (
    <div className="grid gap-2.5">
      <label
        htmlFor={id}
        className="text-caption font-semibold uppercase tracking-caps text-subtle"
      >
        {label}
      </label>

      <div className="group relative">
        <input
          id={id}
          aria-describedby={describedBy}
          aria-invalid={error ? true : undefined}
          className={cn(
            "w-full bg-transparent pb-3 outline-none",
            // 16px+ at every breakpoint, so iOS never zooms on focus.
            "text-[clamp(1.375rem,4vw,2rem)] font-semibold tracking-snug",
            "placeholder:font-normal placeholder:text-subtle",
            className,
          )}
          {...props}
        />

        {/* Track and fill are separate so only transform animates. */}
        <span aria-hidden className="absolute inset-x-0 bottom-0 h-px bg-border-strong" />
        <span
          aria-hidden
          className={cn(
            "absolute inset-x-0 bottom-0 h-[2px] origin-left scale-x-0",
            "transition-transform duration-300 ease-out-soft",
            "group-focus-within:scale-x-100",
          )}
          style={{
            background: error ? "var(--color-danger)" : "var(--lit)",
            boxShadow: error ? "none" : "0 0 12px var(--lit)",
          }}
        />
      </div>

      {error ? (
        <p id={`${id}-error`} role="alert" className="text-small text-danger">
          {error}
        </p>
      ) : hint ? (
        <p id={`${id}-hint`} className="text-small text-subtle">
          {hint}
        </p>
      ) : null}
    </div>
  );
}
