import { Chip, Drawer, DrawerBody, DrawerContent, DrawerTitle, cn } from "@ava/ui";
import type { DeviceDto } from "@ava/contracts";
import { XIcon } from "lucide-react";
import type { ReactNode } from "react";

import { DeviceControls } from "./device-controls";
import { SHEET_SNAP } from "../constants";

/** Whether the device can be reached, and whose fault it is: one value, not two booleans. */
export type Connectivity = "online" | "device-offline" | "hub-offline";

function StatusChip({ connectivity }: { connectivity: Connectivity }) {
  if (connectivity === "online") return null;

  return <Chip tone="warning">{connectivity === "hub-offline" ? "Hub offline" : "Offline"}</Chip>;
}

// Keyed on the device so a half-finished drag cannot carry over into the next one.
function Controls({
  device,
  connectivity,
  onLevelChange,
}: {
  device: DeviceDto;
  connectivity: Connectivity;
  onLevelChange?: (level: number | null) => void;
}) {
  return (
    <DeviceControls
      key={device.id}
      device={device}
      offline={connectivity !== "online"}
      onLevelChange={onLevelChange}
    />
  );
}

type DeviceSheetProps = {
  device: DeviceDto;
  connectivity: Connectivity;
  onClose: () => void;
  onLevelChange?: (level: number | null) => void;
};

/** One device's controls as a column beside the room, for wide screens. */
export function DevicePanel({ device, connectivity, onClose, onLevelChange }: DeviceSheetProps) {
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
            <StatusChip connectivity={connectivity} />
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

      <div className="min-h-0 flex-1 overflow-y-auto px-5 pb-8">
        <Controls device={device} connectivity={connectivity} onLevelChange={onLevelChange} />
      </div>
    </aside>
  );
}

/** One device's controls as a bottom sheet for phones, with `children` as its handle. */
export function DeviceDrawer({
  device,
  connectivity,
  onClose,
  onLevelChange,
  children,
}: DeviceSheetProps & { children: ReactNode }) {
  return (
    <Drawer open onOpenChange={(open) => !open && onClose()} snapPoints={SHEET_SNAP}>
      <DrawerContent className="h-[94dvh] max-h-[94dvh] px-0">
        <DrawerTitle className="sr-only">{device.name}</DrawerTitle>

        <div className="shrink-0 px-5 pt-3">{children}</div>

        <div className="flex shrink-0 items-center gap-2 px-5 pb-4 pt-5">
          <h2 className="min-w-0 truncate text-body font-semibold">{device.name}</h2>
          <StatusChip connectivity={connectivity} />
        </div>

        <DrawerBody className="px-5 pb-8">
          <Controls device={device} connectivity={connectivity} onLevelChange={onLevelChange} />
        </DrawerBody>
      </DrawerContent>
    </Drawer>
  );
}
