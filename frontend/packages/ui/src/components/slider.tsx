import { Slider as SliderPrimitive } from "radix-ui";
import type { ComponentProps } from "react";

import { cn } from "../lib/utils";

const ROOT_CLASS = [
  "relative flex h-[46px] w-full touch-none select-none items-center",
  "rounded-md bg-raised",
  "data-[disabled]:opacity-50",
];

export type SliderProps = ComponentProps<typeof SliderPrimitive.Root> & {
  /** What fills the range: the accent colour, or the light's own colour via `--lit`. */
  tone?: "accent" | "lit";
};

/** The filled slider — the whole surface is the touch target. */
export function Slider({
  className,
  tone = "accent",
  "aria-label": ariaLabel,
  ...props
}: SliderProps) {
  return (
    <SliderPrimitive.Root
      data-slot="slider"
      className={cn(ROOT_CLASS, "overflow-hidden", className)}
      {...props}
    >
      <SliderPrimitive.Track data-slot="track" className="relative h-full w-full">
        <SliderPrimitive.Range
          data-slot="range"
          className={cn(
            "absolute h-full transition-[width] duration-200 ease-out",
            tone === "lit" ? "bg-[var(--lit)]" : "bg-accent",
          )}
        />
      </SliderPrimitive.Track>

      <SliderPrimitive.Thumb
        data-slot="thumb"
        aria-label={ariaLabel}
        className="block size-full outline-none focus-visible:ring-2 focus-visible:ring-fg focus-visible:ring-inset"
      />
    </SliderPrimitive.Root>
  );
}

export type MarkerSliderProps = ComponentProps<typeof SliderPrimitive.Root>;

/** The unfilled slider — a thin thumb marks a position on the track. */
export function MarkerSlider({ className, "aria-label": ariaLabel, ...props }: MarkerSliderProps) {
  return (
    <SliderPrimitive.Root
      data-slot="slider"
      className={cn(ROOT_CLASS, "px-[3px]", className)}
      {...props}
    >
      <SliderPrimitive.Track data-slot="track" className="relative h-full w-full" />

      <SliderPrimitive.Thumb
        data-slot="thumb"
        aria-label={ariaLabel}
        className={cn(
          "block h-7 w-[5px] rounded-full bg-white outline-none",
          "shadow-[0_0_0_1.5px_rgba(0,0,0,0.45)]",
          "focus-visible:shadow-[0_0_0_2px_var(--palette-fg)]",
        )}
      />
    </SliderPrimitive.Root>
  );
}
