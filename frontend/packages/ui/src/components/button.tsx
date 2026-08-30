import { Slot } from "radix-ui";
import { cva, type VariantProps } from "class-variance-authority";
import { Loader2Icon } from "lucide-react";
import type { ComponentProps } from "react";

import { cn } from "../lib/utils";

const button = cva(
  [
    "inline-flex items-center justify-center gap-2 whitespace-nowrap font-semibold",
    "rounded-md",
    "transition-[background-color,border-color,color,opacity,transform] duration-150 ease-out",
    "active:scale-[0.98]",
    "disabled:pointer-events-none disabled:opacity-40",
    "[&_svg]:size-[18px] [&_svg]:shrink-0 [&_svg]:pointer-events-none",
  ],
  {
    variants: {
      variant: {
        primary: "bg-accent text-accent-fg hover:opacity-90",
        secondary: "bg-raised text-fg hover:bg-border",
        ghost: "border border-border-strong text-fg hover:bg-surface",
        quiet: "text-muted hover:text-fg",
        danger: "bg-danger text-white hover:opacity-90",
      },
      size: {
        md: "h-[50px] px-5 text-body",
        /* 44 for a thumb, 36 where there is a pointer. */
        sm: "h-11 px-3 text-small [@media(hover:hover)]:h-9",
        icon: "size-11",
        "icon-sm": "size-9",
      },
      block: { true: "w-full" },
    },
    defaultVariants: { variant: "primary", size: "md" },
  },
);

export type ButtonProps = ComponentProps<"button"> &
  VariantProps<typeof button> & {
    asChild?: boolean;
    loading?: boolean;
  };

export function Button({
  className,
  variant,
  size,
  block,
  asChild,
  loading = false,
  disabled,
  children,
  ...props
}: ButtonProps) {
  const Comp = asChild ? Slot.Root : "button";

  return (
    <Comp
      className={cn(button({ variant, size, block }), className)}
      disabled={disabled || loading}
      aria-busy={loading || undefined}
      data-loading={loading || undefined}
      {...props}
    >
      {loading ? (
        /* The rotation lives on a span wrapper: animating the SVG element
           itself keeps the browser from compositing it on the GPU. */
        <span className="animate-spin [&_svg]:block" aria-hidden>
          <Loader2Icon />
        </span>
      ) : null}
      {children}
    </Comp>
  );
}

export { button as buttonVariants };
