import { Toaster as Sonner } from "sonner";
import type { ComponentProps } from "react";

import { cn } from "../lib/utils";

const WIDTH = "min(400px, calc(100vw - 24px))";

/** Toasts at the top, fully unstyled — Sonner's runtime sheet outranks Tailwind otherwise. */
export function Toaster({ className, toastOptions, ...props }: ComponentProps<typeof Sonner>) {
  return (
    <Sonner
      position="top-center"
      offset={12}
      gap={10}
      visibleToasts={3}
      className={cn("font-sans", className)}
      /* Inline: below 600px Sonner's own full-bleed offsets outrank a stylesheet. */
      style={
        {
          "--width": WIDTH,
          width: WIDTH,
          left: "50%",
          right: "auto",
          transform: "translateX(-50%)",
        } as React.CSSProperties
      }
      toastOptions={{
        ...toastOptions,
        unstyled: true,
        /* Inline for the same reason: below 600px Sonner's sheet sizes toasts off-centre. */
        style: { width: "100%", ...toastOptions?.style },
        classNames: {
          toast: cn(
            /* w-full: sizing the toast independently left it narrower than its container. */
            "pointer-events-auto flex w-full items-center gap-3 rounded-lg p-3.5",
            "border border-border bg-raised text-fg",
            "shadow-[0_20px_44px_-18px_rgb(0_0_0/0.75)]",
          ),
          content: "grid min-w-0 flex-1 gap-1",
          title: "text-small font-semibold leading-snug",
          description: "text-caption leading-relaxed text-muted",
          icon: "grid shrink-0 place-items-center [&>svg]:size-4",
          /* Tone reaches the icon and the edge, never the whole surface. */
          success: "border-success/35 [&_[data-icon]]:text-success",
          error: "border-danger/35 [&_[data-icon]]:text-danger",
          warning: "border-warning/35 [&_[data-icon]]:text-warning",
          info: "[&_[data-icon]]:text-muted",
          loading: "[&_[data-icon]]:text-muted",
          actionButton: cn(
            "shrink-0 rounded-sm bg-accent px-2.5 py-1.5",
            "text-caption font-semibold text-accent-fg",
            "transition-opacity duration-150 ease-out hover:opacity-85",
          ),
          cancelButton: cn(
            "shrink-0 rounded-sm border border-border px-2.5 py-1.5",
            "text-caption font-medium text-muted",
            "transition-colors duration-150 ease-out hover:text-fg",
          ),
          closeButton: cn(
            "order-last grid size-6 shrink-0 place-items-center rounded-full",
            "border border-border bg-surface text-subtle",
            "transition-colors duration-150 ease-out hover:text-fg",
            "[&>svg]:size-3",
          ),
          ...toastOptions?.classNames,
        },
      }}
      {...props}
    />
  );
}
