import { useId, type ReactNode } from "react";

import { cn } from "../lib/utils";
import { Label } from "./label";

type FieldProps = {
  label: string;
  error?: string;
  hint?: string;
  className?: string;
  children: (props: {
    id: string;
    "aria-describedby": string | undefined;
    invalid: boolean;
  }) => ReactNode;
};

function Field({ label, error, hint, className, children }: FieldProps) {
  const id = useId();
  const errorId = `${id}-error`;
  const hintId = `${id}-hint`;
  const describedBy = [error && errorId, hint && hintId].filter(Boolean).join(" ");

  return (
    <div data-slot="field" className={cn("grid gap-2", className)}>
      <Label htmlFor={id}>{label}</Label>

      {children({ id, "aria-describedby": describedBy || undefined, invalid: Boolean(error) })}

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
  );
}

export { Field, type FieldProps };
