import { AlertDialog } from "radix-ui";
import type { ReactNode } from "react";

import { cn } from "../lib/utils";
import { buttonVariants } from "./button";

/**
 * A pause before something that cannot be undone.
 *
 * Radix's AlertDialog rather than a hand-rolled sheet, because the behaviour
 * that matters here is the part nobody writes by hand correctly: focus is
 * trapped and lands on the safe choice, Escape cancels, the page behind goes
 * inert, and the whole thing is announced as an alert rather than as a passing
 * region. A confirmation that can be dismissed by accident is worse than none,
 * because it teaches people the guard is there.
 *
 * Deliberately not used for everything. A scene takes ten seconds to rebuild
 * and asking about it would be nagging; a room takes every device in it down
 * with it, and that is worth a sentence.
 */
export function Confirm({
  open,
  onOpenChange,
  title,
  description,
  confirmLabel,
  onConfirm,
  tone = "default",
  children,
}: {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  title: string;
  /** What will actually happen. Name the consequence, not the action. */
  description: ReactNode;
  confirmLabel: string;
  onConfirm: () => void;
  /** Matches the tone vocabulary of Chip and MenuItem. */
  tone?: "default" | "danger";
  /** The control that opens this, if it should own its own trigger. */
  children?: ReactNode;
}) {
  return (
    <AlertDialog.Root open={open} onOpenChange={onOpenChange}>
      {children ? <AlertDialog.Trigger asChild>{children}</AlertDialog.Trigger> : null}

      <AlertDialog.Portal>
        <AlertDialog.Overlay
          className={cn("fixed inset-0 z-overlay bg-scrim", "data-[state=open]:animate-fade-in")}
        />

        <AlertDialog.Content
          className={cn(
            "fixed left-1/2 top-1/2 z-modal w-[calc(100%-2.5rem)] max-w-[380px]",
            "-translate-x-1/2 -translate-y-1/2 rounded-2xl border border-border bg-surface",
            "p-5 dialog-lift outline-none",
            "data-[state=open]:animate-fade-in",
          )}
        >
          <AlertDialog.Title className="text-title font-semibold">{title}</AlertDialog.Title>

          <AlertDialog.Description className="mt-2 text-small leading-relaxed text-muted">
            {description}
          </AlertDialog.Description>

          <div className="mt-5 grid grid-cols-2 gap-2.5">
            <AlertDialog.Cancel className={cn(buttonVariants({ variant: "secondary" }))}>
              Cancel
            </AlertDialog.Cancel>

            <AlertDialog.Action
              onClick={onConfirm}
              className={cn(buttonVariants({ variant: tone === "danger" ? "danger" : "primary" }))}
            >
              {confirmLabel}
            </AlertDialog.Action>
          </div>
        </AlertDialog.Content>
      </AlertDialog.Portal>
    </AlertDialog.Root>
  );
}
