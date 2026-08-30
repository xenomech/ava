import { use, type ComponentProps } from "react";

import { cn } from "../lib/utils";
import { FieldContext } from "./field";

type InputProps = ComponentProps<"input"> & {
  invalid?: boolean;
};

function Input({
  className,
  invalid,
  id,
  "aria-describedby": ariaDescribedBy,
  type = "text",
  ...props
}: InputProps) {
  /* Inside a <Field>, wiring comes from context; explicit props still win. */
  const field = use(FieldContext);
  const isInvalid = invalid ?? field?.invalid;

  return (
    <input
      type={type}
      id={id ?? field?.id}
      data-slot="input"
      className={cn(
        "h-[52px] w-full rounded-md border bg-surface px-4",
        "text-[16px] placeholder:text-subtle",
        "transition-colors duration-150 ease-out",
        "focus:border-fg focus:outline-none",
        "disabled:cursor-not-allowed disabled:opacity-50",
        isInvalid && "border-danger focus:border-danger",
        className,
      )}
      aria-invalid={isInvalid || undefined}
      aria-describedby={ariaDescribedBy ?? field?.describedBy}
      {...props}
    />
  );
}

export { Input, type InputProps };
