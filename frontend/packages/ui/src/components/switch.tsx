import { Switch as SwitchPrimitive } from "radix-ui";
import type { ComponentProps } from "react";

import { cn } from "../lib/utils";

export function Switch({ className, ...props }: ComponentProps<typeof SwitchPrimitive.Root>) {
  return (
    <SwitchPrimitive.Root
      className={cn(
        "peer inline-flex h-7 w-[46px] shrink-0 items-center rounded-full p-[3px]",
        "bg-off data-[state=checked]:bg-accent",
        "transition-colors duration-200 ease-out",
        "disabled:cursor-not-allowed disabled:opacity-50",
        className,
      )}
      {...props}
    >
      <SwitchPrimitive.Thumb
        className={cn(
          "pointer-events-none block size-[22px] rounded-full",
          "bg-surface data-[state=checked]:bg-accent-fg",
          "translate-x-0 data-[state=checked]:translate-x-[18px]",
          "transition-transform duration-[380ms] ease-spring",
        )}
      />
    </SwitchPrimitive.Root>
  );
}
