import { createContext, use, useId, useMemo, type ReactNode } from "react";

import { Slot } from "radix-ui";

import { cn } from "../lib/utils";
import { Label } from "./label";

type FieldContextValue = {
  id: string;
  describedBy: string | undefined;
  invalid: boolean;
};

// Controls inside a Field pick up their id, description, and validity from here.
const FieldContext = createContext<FieldContextValue | null>(null);

function useField() {
  const field = use(FieldContext);
  if (!field) throw new Error("useField must be used inside <Field>");
  return field;
}

type FieldProps = {
  label: string;
  error?: string;
  hint?: string;
  /** Sits on the label row, right-aligned — a "Forgot?" link, a unit toggle. */
  action?: ReactNode;
  className?: string;
  children: ReactNode;
};

function Field({ label, error, hint, action, className, children }: FieldProps) {
  const id = useId();
  const errorId = `${id}-error`;
  const hintId = `${id}-hint`;
  // The error replaces the hint, so only point at the one actually in the document.
  const describedBy = error ? errorId : hint ? hintId : undefined;
  const invalid = Boolean(error);

  const field = useMemo(() => ({ id, describedBy, invalid }), [id, describedBy, invalid]);

  return (
    <FieldContext value={field}>
      <div data-slot="field" className={cn("grid gap-2", className)}>
        {action ? (
          <div className="flex items-baseline justify-between gap-3">
            <Label htmlFor={id}>{label}</Label>
            {action}
          </div>
        ) : (
          <Label htmlFor={id}>{label}</Label>
        )}

        {children}

        {hint && !error ? (
          <p id={hintId} className="text-caption text-subtle">
            {hint}
          </p>
        ) : null}

        {error ? (
          <p id={errorId} role="alert" className="text-caption text-danger">
            {error}
          </p>
        ) : null}
      </div>
    </FieldContext>
  );
}

/** Marks its child as the field's control, for elements that don't read FieldContext. */
function FieldControl({ children }: { children: ReactNode }) {
  const { id, describedBy, invalid } = useField();
  return (
    <Slot.Root id={id} aria-describedby={describedBy} aria-invalid={invalid || undefined}>
      {children}
    </Slot.Root>
  );
}

export { Field, FieldContext, FieldControl, useField, type FieldProps };
