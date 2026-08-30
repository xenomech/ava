import { Device, cn } from "@ava/ui";
import type { DeviceDto } from "@ava/contracts";
import { Link } from "@tanstack/react-router";

import { AddDevices } from "./add-devices";
import { deviceColor, deviceKind, deviceLabel, deviceLevel } from "../lib/device-view";
import { OnAPlug } from "./on-a-plug";

/**
 * The room's devices as a carousel you can thumb through.
 *
 * Every card is the same width and shows the device itself, lit at whatever it
 * is actually doing, so the row is scannable at arm's length without reading a
 * word. Names truncate rather than wrap: a card that grows to fit its label
 * breaks the rhythm of the row and lands the next card somewhere different
 * every time.
 *
 * The scrollbar is hidden. A card left half off the edge already says the row
 * scrolls, and that is the only affordance a touch screen was ever going to
 * get.
 *
 * On a phone this row is also the sheet's handle: selecting a device lifts the
 * strip and unfolds the controls beneath it, so browsing devices and dismissing
 * the sheet are the same object. That is why the cards are links rather than
 * buttons — the selection lives in the URL, so the back button closes it and a
 * device stays addressable on its own.
 */
export function DeviceStrip({
  devices,
  elsewhere,
  roomId,
  roomName,
  selectedId,
  label,
  className,
}: {
  devices: DeviceDto[];
  /** Everything not in this room, offered by the trailing add card. */
  elsewhere: DeviceDto[];
  roomId: string;
  roomName: string;
  selectedId?: string;
  label: string;
  className?: string;
}) {
  return (
    <div
      aria-label={label}
      /* Vaul reads a vertical drag anywhere in the sheet as a dismiss. Without
         this the row cannot be scrolled sideways without the sheet following. */
      data-vaul-no-drag
      className={cn(
        "no-scrollbar -mx-5 flex snap-x snap-mandatory gap-2.5 overflow-x-auto px-5 pb-1",
        /* scroll-padding, not padding: opening the sheet scrolls the selected
           card into view, and without this it lands flush against the edge
           with the padding scrolled out of sight. */
        "scroll-pl-5 sm:-mx-6 sm:scroll-pl-6 sm:px-6",
        className,
      )}
    >
      {devices.map((device) => (
        <DeviceCard
          key={device.id}
          device={device}
          roomId={roomId}
          selected={device.id === selectedId}
        />
      ))}

      {/* Last in the row rather than a button in the header: adding a device is
          the same kind of act as picking one, and it belongs where the room's
          devices already are. */}
      <AddDevices roomId={roomId} roomName={roomName} candidates={elsewhere} />
    </div>
  );
}

function DeviceCard({
  device,
  roomId,
  selected,
}: {
  device: DeviceDto;
  roomId: string;
  selected: boolean;
}) {
  const level = deviceLevel(device);
  const color = deviceColor(device);

  return (
    <Link
      to="/rooms/$roomId"
      params={{ roomId }}
      search={selected ? {} : { device: device.id }}
      replace={selected}
      aria-current={selected ? "true" : undefined}
      style={{ "--level": level, "--lit": color } as React.CSSProperties}
      className={cn(
        "relative grid w-[126px] shrink-0 snap-start content-start gap-1.5 rounded-lg p-2.5",
        "border bg-surface transition-colors duration-150 ease-out",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
        selected ? "border-fg" : "border-border hover:border-border-strong",
        device.status === "offline" && "opacity-55",
      )}
    >
      {device.kind === "plug" && device.appliance ? (
        <OnAPlug className="absolute right-2 top-2 size-5" />
      ) : null}

      <span className="grid h-14 w-full place-items-center">
        <Device kind={deviceKind(device)} level={level} color={color} className="h-full" />
      </span>

      <span className="w-full truncate text-small font-semibold">{device.name}</span>
      <span className="w-full truncate font-mono text-caption text-subtle tabular">
        {deviceLabel(device)}
      </span>
    </Link>
  );
}
