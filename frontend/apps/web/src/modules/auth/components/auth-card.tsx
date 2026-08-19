import type { ReactNode } from "react";

export function AuthCard({
  title,
  description,
  children,
  footer,
}: {
  title: string;
  description?: string;
  children?: ReactNode;
  footer?: ReactNode;
}) {
  return (
    <main className="grid min-h-full place-items-center bg-bg p-6">
      <div className="w-full max-w-[400px] rounded-xl border border-border bg-surface p-7">
        <div className="mb-6 flex items-center gap-2.5">
          <span className="grid size-7 place-items-center rounded-xs bg-accent text-small font-bold text-accent-fg">
            a
          </span>
          <span className="text-lead font-semibold">ava</span>
        </div>

        <header className="mb-5 grid gap-1.5">
          <h1 className="text-display font-semibold text-balance">{title}</h1>
          {description ? <p className="text-small text-muted">{description}</p> : null}
        </header>

        {children ? <div className="grid gap-4">{children}</div> : null}

        {footer ? (
          <div className="mt-6 border-t border-border pt-4 text-small text-muted">{footer}</div>
        ) : null}
      </div>
    </main>
  );
}
