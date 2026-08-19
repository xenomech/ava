import { Slider as SliderPrimitive } from "radix-ui";
import type { ComponentProps } from "react";

import { cn } from "../lib/utils";

export type SliderProps = ComponentProps<typeof SliderPrimitive.Root> & {
  lit?: boolean;
  variant?: "fill" | "marker";
};

export function Slider({ className, lit, variant = "fill", ...props }: SliderProps) {
  const marker = variant === "marker";

  return (
    <SliderPrimitive.Root
      data-slot="slider"
      className={cn(
        "relative flex h-[46px] w-full touch-none select-none items-center",
        "rounded-md bg-raised overflow-hidden",
        "data-[disabled]:opacity-50",
        className,
      )}
      {...props}
    >
      <SliderPrimitive.Track data-slot="track" className="relative h-full w-full">
        <SliderPrimitive.Range
          data-slot="range"
          className={cn(
            "absolute h-full transition-[width] duration-200 ease-out",
            marker ? "bg-transparent" : lit ? "bg-[var(--lit)]" : "bg-accent",
          )}
        />
      </SliderPrimitive.Track>

      <SliderPrimitive.Thumb
        data-slot="thumb"
        aria-label={props["aria-label"]}
        className={cn(
          "block outline-none",
          marker
            ? [
                "h-7 w-[5px] rounded-full bg-white",
                "shadow-[0_0_0_1.5px_rgba(0,0,0,0.45)]",
                "focus-visible:shadow-[0_0_0_2px_var(--palette-fg)]",
              ]
            : "size-full focus-visible:ring-2 focus-visible:ring-fg focus-visible:ring-inset",
        )}
      />
    </SliderPrimitive.Root>
  );
}
