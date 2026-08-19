import { Drawer as Vaul } from "vaul";
import type { ComponentProps, ReactNode } from "react";

import { cn } from "../lib/utils";

function Drawer(props: ComponentProps<typeof Vaul.Root>) {
  return <Vaul.Root setBackgroundColorOnScale={false} {...props} />;
}

const DrawerTrigger = Vaul.Trigger;
const DrawerClose = Vaul.Close;

export {
  Drawer,
  DrawerTrigger,
  DrawerClose,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerDescription,
  DrawerBody,
  DrawerFooter,
};

function DrawerContent({
  className,
  children,
  ...props
}: ComponentProps<typeof Vaul.Content> & { children: ReactNode }) {
  return (
    <Vaul.Portal>
      <Vaul.Overlay className="fixed inset-0 z-overlay bg-scrim" />
      <Vaul.Content
        className={cn(
          "fixed inset-x-0 bottom-0 z-modal mx-auto flex max-h-[92dvh] flex-col",
          "rounded-t-2xl border-t border-border bg-surface outline-none",
          "w-full sm:max-w-[520px]",
          className,
        )}
        {...props}
      >
        <div
          aria-hidden
          className="mx-auto mt-3 h-[5px] w-10 shrink-0 rounded-full bg-border-strong"
        />
        {children}
      </Vaul.Content>
    </Vaul.Portal>
  );
}

function DrawerHeader({ className, ...props }: ComponentProps<"header">) {
  return (
    <header
      className={cn("flex shrink-0 items-start justify-between gap-4 px-5 pb-3 pt-4", className)}
      {...props}
    />
  );
}

function DrawerTitle({ className, ...props }: ComponentProps<typeof Vaul.Title>) {
  return (
    <Vaul.Title
      className={cn("text-title font-semibold", className)}
      {...props}
    />
  );
}

function DrawerDescription({ className, ...props }: ComponentProps<typeof Vaul.Description>) {
  return <Vaul.Description className={cn("mt-1 text-small text-muted", className)} {...props} />;
}

function DrawerBody({ className, ...props }: ComponentProps<"div">) {
  return (
    <div
      className={cn("min-h-0 flex-1 overflow-y-auto overscroll-contain px-5", className)}
      {...props}
    />
  );
}

function DrawerFooter({ className, ...props }: ComponentProps<"footer">) {
  return (
    <footer
      className={cn(
        "flex shrink-0 gap-2.5 border-t border-border bg-surface px-5 pt-3.5",
        "pb-[calc(1.125rem+env(safe-area-inset-bottom))]",
        className,
      )}
      {...props}
    />
  );
}
