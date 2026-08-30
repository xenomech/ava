import { DropdownMenu as MenuPrimitive } from "radix-ui";
import type { ComponentProps } from "react";

import { cn } from "../lib/utils";

function Menu(props: ComponentProps<typeof MenuPrimitive.Root>) {
  return <MenuPrimitive.Root {...props} />;
}

const MenuTrigger = MenuPrimitive.Trigger;

function MenuContent({
  className,
  align = "start",
  sideOffset = 6,
  ...props
}: ComponentProps<typeof MenuPrimitive.Content>) {
  return (
    <MenuPrimitive.Portal>
      <MenuPrimitive.Content
        align={align}
        sideOffset={sideOffset}
        className={cn(
          "z-modal min-w-[200px] overflow-hidden rounded-md border border-border bg-raised p-1",
          "shadow-[0_20px_44px_-18px_rgb(0_0_0/0.75)]",
          "animate-scale-in origin-(--radix-dropdown-menu-content-transform-origin)",
          "motion-reduce:animate-none",
          className,
        )}
        {...props}
      />
    </MenuPrimitive.Portal>
  );
}

/** `tone="danger"` for anything destructive; 44 tall on touch, 40 with a pointer. */
function MenuItem({
  className,
  tone,
  ...props
}: ComponentProps<typeof MenuPrimitive.Item> & { tone?: "danger" }) {
  return (
    <MenuPrimitive.Item
      className={cn(
        "flex h-11 cursor-default items-center gap-3 rounded-xs px-3 text-body outline-none select-none",
        "[@media(hover:hover)]:h-10",
        "[&_svg]:size-4 [&_svg]:shrink-0",
        "data-[highlighted]:bg-surface",
        "data-[disabled]:pointer-events-none data-[disabled]:opacity-40",
        tone === "danger" ? "text-danger" : "text-fg",
        className,
      )}
      {...props}
    />
  );
}

function MenuSeparator({ className, ...props }: ComponentProps<typeof MenuPrimitive.Separator>) {
  return <MenuPrimitive.Separator className={cn("my-1 h-px bg-border", className)} {...props} />;
}

export { Menu, MenuTrigger, MenuContent, MenuItem, MenuSeparator };
