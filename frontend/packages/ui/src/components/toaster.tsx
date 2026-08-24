import { Toaster as Sonner } from "sonner";
import type { ComponentProps } from "react";

import { cn } from "../lib/utils";

const WIDTH = "min(400px, calc(100vw - 24px))";

/**
 * Toasts, in the app's own language.
 *
 * Three deliberate departures from Sonner's defaults.
 *
 * They sit at the top rather than the bottom, because the bottom of every
 * screen here is already spoken for — the device strip on a room, the sheet
 * that rises over it, the controls along the foot of a device. A toast down
 * there lands on the thing the person is using.
 *
 * `richColors` is off. It paints the whole toast a stock green or red that
 * belongs to no part of this app; tone is carried the way `Chip` carries it, by
 * the icon and a tinted edge over the ordinary raised surface.
 *
 * And everything is `unstyled`. Sonner injects its stylesheet at runtime, so it
 * lands after Tailwind and wins every specificity tie — a `bg-raised` on the
 * toast simply loses to `[data-sonner-toast]`. Its defaults are all guarded
 * behind `data-styled=true`, which `unstyled` turns off, so this is the
 * supported way out rather than an escalation war.
 */
export function Toaster({ className, toastOptions, ...props }: ComponentProps<typeof Sonner>) {
  return (
    <Sonner
      position="top-center"
      offset={12}
      gap={10}
      visibleToasts={3}
      className={cn("font-sans", className)}
      /* Inline, because below 600px Sonner swaps top-center for its own
         full-bleed left/right offsets and the container ends up wider than the
         screen. A stylesheet cannot outrank that reliably; an inline style can,
         and it keeps one centred column at every width. */
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
        classNames: {
          toast: cn(
            "pointer-events-auto flex w-[var(--width)] items-center gap-3 rounded-lg p-3.5",
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
