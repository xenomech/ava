import { Drawer, DrawerContent, DrawerDescription, DrawerTitle, cn } from "@ava/ui";
import { useRef, useState, type PointerEvent, type ReactNode } from "react";

/** A downward drag this far on the handle counts as pulling the menu open. */
const PULL = 10;

/**
 * The navigation, pulled down from the top edge.
 *
 * A hamburger in the corner is the wrong shape for a phone held in one hand:
 * it is the furthest point from the thumb, and it spends the whole time
 * floating over content it has nothing to do with. This is the gesture people
 * already know from the status bar — a handle at the top edge that you tap, or
 * drag down, and the panel comes with it.
 *
 * Drag-to-open is wired here rather than left to the sheet, which only knows
 * how to be dragged once it is already open.
 */
export function MenuHandle({ children }: { children: ReactNode }) {
  const [open, setOpen] = useState(false);
  const from = useRef<number | null>(null);

  const down = (event: PointerEvent<HTMLButtonElement>) => {
    from.current = event.clientY;
  };

  const move = (event: PointerEvent<HTMLButtonElement>) => {
    if (from.current === null || open) return;

    if (event.clientY - from.current > PULL) {
      from.current = null;
      setOpen(true);
    }
  };

  const up = () => {
    from.current = null;
  };

  return (
    <>
      <button
        type="button"
        aria-label="Open menu"
        aria-expanded={open}
        onClick={() => setOpen(true)}
        onPointerDown={down}
        onPointerMove={move}
        onPointerUp={up}
        onPointerCancel={up}
        className={cn(
          /* A wide, shallow target: easy to hit without looking, and it never
             covers anything because the page starts below it. */
          "group absolute inset-x-0 top-0 z-sticky grid h-11 touch-none place-items-center md:hidden",
          "focus-visible:outline-none",
        )}
      >
        <span
          aria-hidden
          className={cn(
            "h-1 w-10 rounded-full bg-border-strong",
            "transition-[background-color,width] duration-200 ease-out",
            "group-hover:w-12 group-hover:bg-muted group-focus-visible:bg-fg",
          )}
        />
      </button>

      <Drawer open={open} onOpenChange={setOpen} direction="top">
        <DrawerContent
          className={cn(
            "inset-x-0 bottom-auto top-0 mt-0 max-h-[88dvh] rounded-t-none rounded-b-2xl",
            "border-b border-t-0 pt-[max(0.5rem,env(safe-area-inset-top))]",
          )}
          grabber={false}
        >
          <DrawerTitle className="sr-only">Menu</DrawerTitle>
          <DrawerDescription className="sr-only">Rooms and account</DrawerDescription>

          <div className="min-h-0 flex-1 overflow-y-auto px-4 pb-4 pt-2">{children}</div>

          {/* The handle repeats at the bottom of the sheet, where the panel now
              ends and the thumb already is. */}
          <span aria-hidden className="mx-auto mb-3 h-1 w-10 shrink-0 rounded-full bg-border-strong" />
        </DrawerContent>
      </Drawer>
    </>
  );
}
