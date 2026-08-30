import { Chip, Drawer, DrawerBody, DrawerContent, DrawerTitle, cn } from "@ava/ui";
import type { DeviceDto } from "@ava/contracts";
import { XIcon } from "lucide-react";
import type { ReactNode } from "react";

import { DeviceControls } from "./device-controls";

/**
 * Where the phone sheet rests before you drag it.
 *
 * Half the screen, so the device is still visible on the stage above it — the
 * whole point of swapping the room's switch for the device is lost if the sheet
 * covers it. Power and brightness fit in that half; colour and details are the
 * reward for pulling it further. `ROOM_HEIGHT` is the other side of the same
 * number, and the room caps its own height to match.
 */
export const SHEET_SNAP = [0.52, 0.94];
export const ROOM_HEIGHT = "h-[48dvh]";

/**
 * One device's controls, at whichever edge suits the screen.
 *
 * On a phone it is a bottom sheet whose handle is the room's own device strip,
 * so browsing devices and dismissing the sheet are the same object rather than
 * two things competing for the bottom of the screen. On a desktop the strip
 * stays in the room's footer and this becomes a column beside the switch — a
 * short wide band across a 1440px screen wastes the width and puts the controls
 * a long way from the pointer.
 *
 * The desktop side is deliberately not a drawer. It never needs to be dragged
 * away, so it is a panel with a transition rather than a gesture surface.
 */
export function DeviceSheet({
  device,
  offline,
  hubOffline,
  strip,
  beside,
  onClose,
  onLevelChange,
}: {
  device: DeviceDto;
  offline: boolean;
  hubOffline: boolean;
  /** The room's device strip, used as the phone sheet's handle. */
  strip: ReactNode;
  /** Render as a column beside the room rather than a sheet under it. */
  beside: boolean;
  onClose: () => void;
  onLevelChange?: (level: number | null) => void;
}) {
  const status = offline ? (
    <Chip tone="warning">{hubOffline ? "Hub offline" : "Offline"}</Chip>
  ) : null;

  /* Keyed on the device so a half-finished slider drag cannot carry over into
     whichever device is picked next. */
  const controls = (
    <DeviceControls
      key={device.id}
      device={device}
      offline={offline}
      onLevelChange={onLevelChange}
    />
  );

  if (beside) {
    return (
      <aside
        aria-label={`${device.name} controls`}
        className={cn(
          "flex w-[340px] min-h-0 flex-col border-l border-border bg-surface",
          "animate-slide-in-end",
        )}
      >
        <header className="flex shrink-0 items-start justify-between gap-3 p-5 pb-4">
          <div className="min-w-0">
            <h2 className="truncate text-title font-semibold">{device.name}</h2>
            <p className="mt-1 flex items-center gap-2 font-mono text-caption text-subtle">
              {device.room || "No room"}
              {status}
            </p>
          </div>

          <button
            type="button"
            onClick={onClose}
            aria-label="Close controls"
            className={cn(
              "grid size-9 shrink-0 place-items-center rounded-full text-muted",
              "transition-colors duration-150 ease-out hover:bg-raised hover:text-fg",
              "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
            )}
          >
            <XIcon className="size-4" aria-hidden />
          </button>
        </header>

        <div className="min-h-0 flex-1 overflow-y-auto px-5 pb-8">{controls}</div>
      </aside>
    );
  }

  return (
    <Drawer open onOpenChange={(open) => !open && onClose()} snapPoints={SHEET_SNAP}>
      <DrawerContent className="h-[94dvh] max-h-[94dvh] px-0">
        <DrawerTitle className="sr-only">{device.name}</DrawerTitle>

        <div className="shrink-0 px-5 pt-3">{strip}</div>

        <div className="flex shrink-0 items-center gap-2 px-5 pb-4 pt-5">
          <h2 className="min-w-0 truncate text-body font-semibold">{device.name}</h2>
          {status}
        </div>

        <DrawerBody className="px-5 pb-8">{controls}</DrawerBody>
      </DrawerContent>
    </Drawer>
  );
}
