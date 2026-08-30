import { Button, cn } from "@ava/ui";
import type { ReactNode } from "react";

import { Reveal } from "@/shared/components/reveal";

// One question, filling the screen. No card, no border, nothing competing with
// the sentence you are being asked. Everything cascades in from the top down so
// the eye is led through it in reading order rather than arriving at once.
export function StepFrame({
  eyebrow,
  title,
  description,
  children,
  error,
  action,
  onSubmit,
}: {
  eyebrow: ReactNode;
  title: string;
  description?: string;
  children: ReactNode;
  error?: string | null;
  /**
   * The single way forward. There is deliberately never a second one: on the
   * hub step an un-paired Continue could only ever fail, so instead of sitting
   * next to Skip looking like its twin, it is simply not rendered.
   */
  action: {
    label: string;
    loading?: boolean;
    variant?: "primary" | "ghost";
    /** Given, the button is an ordinary button. Omitted, it submits the form. */
    onClick?: () => void;
  };
  onSubmit: () => void;
}) {
  return (
    <form
      className="grid gap-9"
      onSubmit={(event) => {
        event.preventDefault();
        onSubmit();
      }}
    >
      <div className="grid gap-4">
        <Reveal at={0}>{eyebrow}</Reveal>

        <Reveal at={1}>
          <h1
            className={cn(
              "text-balance font-semibold tracking-tighter",
              "text-[clamp(2.25rem,7vw,3.75rem)] leading-[1.02]",
            )}
          >
            {title}
          </h1>
        </Reveal>

        {description ? (
          <Reveal at={2}>
            <p className="max-w-[46ch] text-pretty text-lead text-muted">{description}</p>
          </Reveal>
        ) : null}
      </div>

      <Reveal at={3}>{children}</Reveal>

      {error ? (
        <Reveal at={4}>
          <p role="alert" className="text-small text-danger">
            {error}
          </p>
        </Reveal>
      ) : null}

      <Reveal at={5}>
        <Button
          type={action.onClick ? "button" : "submit"}
          size="md"
          variant={action.variant ?? "primary"}
          loading={action.loading}
          onClick={action.onClick}
          className="min-w-[168px]"
        >
          {action.label}
        </Button>
      </Reveal>
    </form>
  );
}

// Where you are, without repeating the headline underneath it. The segments
// carry the progress; the word stays constant so the eye ignores it after the
// first screen.
export function Eyebrow({ at, total }: { at: number; total: number }) {
  return (
    <div className="flex items-center gap-3.5">
      <span className="text-caption font-semibold uppercase tracking-caps text-subtle">Setup</span>

      <span aria-hidden className="flex items-center gap-1.5">
        {Array.from({ length: total }, (_, index) => (
          <span
            key={index}
            className="h-[3px] w-7 rounded-full transition-colors duration-500 ease-out-soft"
            style={{
              background: index <= at ? "var(--lit)" : "var(--color-border-strong)",
              boxShadow: index <= at ? "0 0 10px -2px var(--lit)" : "none",
            }}
          />
        ))}
      </span>

      <span className="text-caption text-subtle tabular">
        {at + 1} of {total}
      </span>
    </div>
  );
}
