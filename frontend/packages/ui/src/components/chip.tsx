import { Slot, Toggle as TogglePrimitive } from "radix-ui";
import { cva, type VariantProps } from "class-variance-authority";
import type { ComponentProps } from "react";

import { cn } from "../lib/utils";

const chip = cva(
  [
    "inline-flex items-center gap-2 whitespace-nowrap rounded-sm border px-3.5 py-2",
    "text-small font-medium",
    "transition-colors duration-150 ease-out",
    "[&_svg]:size-3.5 [&_svg]:shrink-0",
  ],
  {
    variants: {
      tone: {
        neutral: "border-border bg-surface text-fg",
        muted: "border-border bg-transparent text-muted",
        success: "border-success/40 text-success",
        warning: "border-warning/40 text-warning",
        danger: "border-danger/40 text-danger",
      },
      selectable: {
        true: "data-[state=on]:border-accent data-[state=on]:bg-accent data-[state=on]:text-accent-fg",
      },
    },
    defaultVariants: { tone: "neutral" },
  },
);

/* `selectable` is ChipToggle's half of the cva — exposing it here would let it
   leak onto the <span> as an unknown DOM attribute. */
export type ChipProps = ComponentProps<"span"> &
  Omit<VariantProps<typeof chip>, "selectable"> & { asChild?: boolean };

export function Chip({ className, tone, asChild, ...props }: ChipProps) {
  const Comp = asChild ? Slot.Root : "span";
  return <Comp className={cn(chip({ tone }), className)} {...props} />;
}

export function ChipToggle({
  className,
  tone,
  ...props
}: ComponentProps<typeof TogglePrimitive.Root> & VariantProps<typeof chip>) {
  return (
    <TogglePrimitive.Root className={cn(chip({ tone, selectable: true }), className)} {...props} />
  );
}
