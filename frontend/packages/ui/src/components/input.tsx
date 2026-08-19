import type { ComponentProps } from "react";

import { cn } from "../lib/utils";

type InputProps = ComponentProps<"input"> & {
  invalid?: boolean;
};

function Input({ className, invalid, type = "text", ...props }: InputProps) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        "h-[52px] w-full rounded-md border bg-surface px-4",
        "text-[16px] placeholder:text-subtle",
        "transition-colors duration-150 ease-out",
        "focus:border-fg focus:outline-none",
        "disabled:cursor-not-allowed disabled:opacity-50",
        invalid && "border-danger focus:border-danger",
        className,
      )}
      aria-invalid={invalid || undefined}
      {...props}
    />
  );
}

export { Input, type InputProps };
