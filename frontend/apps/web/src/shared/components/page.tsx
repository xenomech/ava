import { cn } from "@ava/ui";
import type { ReactNode } from "react";

export function Page({
  title,
  description,
  actions,
  children,
  className,
}: {
  /** Omitted when a parent layout already owns the heading, as Settings does. */
  title?: string;
  description?: string;
  actions?: ReactNode;
  children: ReactNode;
  className?: string;
}) {
  return (
    <div className={cn("mx-auto grid max-w-[720px] gap-6 px-5 py-8 sm:px-8", className)}>
      {title || description || actions ? (
        <header className="flex items-start justify-between gap-4">
          <div className="grid gap-1">
            {title ? <h1 className="text-title font-semibold">{title}</h1> : null}
            {description ? <p className="text-small text-muted">{description}</p> : null}
          </div>
          {actions}
        </header>
      ) : null}

      {children}
    </div>
  );
}

export function Section({
  title,
  description,
  children,
  className,
}: {
  title?: string;
  description?: string;
  children: ReactNode;
  className?: string;
}) {
  return (
    <section className={cn("rounded-lg border border-border bg-surface", className)}>
      {title ? (
        <div className="grid gap-1 border-b border-border px-5 py-4">
          <h2 className="text-small font-semibold">{title}</h2>
          {description ? <p className="text-caption text-muted">{description}</p> : null}
        </div>
      ) : null}

      {children}
    </section>
  );
}
