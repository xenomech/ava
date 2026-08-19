import { Command as CommandPrimitive } from "cmdk";
import { Dialog as DialogPrimitive } from "radix-ui";
import type { ComponentProps, ReactNode } from "react";

import { cn } from "../lib/utils";

function CommandDialog({
  open,
  onOpenChange,
  label = "Command palette",
  children,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  label?: string;
  children: ReactNode;
}) {
  return (
    <DialogPrimitive.Root open={open} onOpenChange={onOpenChange}>
      <DialogPrimitive.Portal>
        <DialogPrimitive.Overlay className="fixed inset-0 z-overlay bg-scrim data-[state=open]:animate-fade-in motion-reduce:animate-none" />

        <DialogPrimitive.Content
          aria-label={label}
          className={cn(
            "fixed top-[12vh] left-1/2 z-modal w-[calc(100vw-2rem)] max-w-[520px] -translate-x-1/2",
            "overflow-hidden rounded-lg border border-border bg-surface",
            "data-[state=open]:animate-scale-in motion-reduce:animate-none",
          )}
        >
          <DialogPrimitive.Title className="sr-only">{label}</DialogPrimitive.Title>
          <CommandPrimitive className="grid">{children}</CommandPrimitive>
        </DialogPrimitive.Content>
      </DialogPrimitive.Portal>
    </DialogPrimitive.Root>
  );
}

function CommandInput({ className, ...props }: ComponentProps<typeof CommandPrimitive.Input>) {
  return (
    <CommandPrimitive.Input
      className={cn(
        "h-14 border-b border-border bg-transparent px-4",
        "text-[16px] text-fg placeholder:text-subtle",
        "outline-none",
        className,
      )}
      {...props}
    />
  );
}

function CommandList({ className, ...props }: ComponentProps<typeof CommandPrimitive.List>) {
  return (
    <CommandPrimitive.List
      className={cn("max-h-[min(52vh,380px)] overflow-y-auto overscroll-contain p-2", className)}
      {...props}
    />
  );
}

function CommandEmpty({ className, ...props }: ComponentProps<typeof CommandPrimitive.Empty>) {
  return (
    <CommandPrimitive.Empty
      className={cn("px-3 py-8 text-center text-small text-muted", className)}
      {...props}
    />
  );
}

function CommandItem({ className, ...props }: ComponentProps<typeof CommandPrimitive.Item>) {
  return (
    <CommandPrimitive.Item
      className={cn(
        "flex h-11 cursor-default items-center justify-between gap-3 rounded-sm px-3 text-body",
        "transition-colors duration-100 ease-out",
        "data-[selected=true]:bg-raised data-[selected=true]:text-fg",
        className,
      )}
      {...props}
    />
  );
}

function CommandHint({ className, ...props }: ComponentProps<"span">) {
  return <span className={cn("shrink-0 text-caption text-subtle", className)} {...props} />;
}

const CommandGroup = CommandPrimitive.Group;

export {
  CommandDialog,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandHint,
};
