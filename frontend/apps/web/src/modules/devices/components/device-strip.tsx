import { cn } from "@ava/ui";
import { isOn, type DeviceDto } from "@ava/contracts";
import { Link } from "@tanstack/react-router";

import { OnAPlug, deviceColor, deviceLabel } from "./device-stage";

/**
 * The room's devices as a row you can thumb through.
 *
 * On a phone this row is also the sheet's handle: selecting a device lifts the
 * strip and unfolds the controls beneath it, so the way you browse devices and
 * the way you close the sheet are the same object. That is why the cards are
 * links rather than buttons — the selection lives in the URL, so the back
 * button closes the sheet and a device is still addressable on its own.
 */
export function DeviceStrip({
  devices,
  roomId,
  selectedId,
  label,
  className,
}: {
  devices: DeviceDto[];
  roomId: string;
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
        /* scroll-padding, not padding: opening the sheet scrolls the selected
           card into view, and without this it lands flush against the edge
           with the padding scrolled out of sight. */
        "-mx-5 flex snap-x snap-mandatory gap-2.5 overflow-x-auto px-5 pb-1",
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
    </div>
  );
}

/** A device reduced to what you need to pick it out: a lit dot and a name. */
function DeviceCard({
  device,
  roomId,
  selected,
}: {
  device: DeviceDto;
  roomId: string;
  selected: boolean;
}) {
  const live = isOn(device) && device.status !== "offline";

  return (
    <Link
      to="/rooms/$roomId"
      params={{ roomId }}
      search={selected ? {} : { device: device.id }}
      replace={selected}
      aria-current={selected ? "true" : undefined}
      className={cn(
        "relative grid w-[168px] shrink-0 snap-start content-start gap-1.5 rounded-lg p-3",
        "border bg-surface transition-colors duration-150 ease-out",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
        selected ? "border-fg" : "border-border hover:border-border-strong",
        device.status === "offline" && "opacity-55",
      )}
    >
      {device.kind === "plug" && device.appliance ? (
        <OnAPlug className="absolute right-2 top-2 size-5" />
      ) : null}

      <span className="flex items-center gap-2 pr-6">
        <span
          aria-hidden
          className={cn("size-1.5 shrink-0 rounded-full", !live && "bg-border-strong")}
          style={live ? { background: deviceColor(device) } : undefined}
        />
        <span className="truncate text-small font-semibold">{device.name}</span>
      </span>
      <span className="font-mono text-caption text-subtle tabular">{deviceLabel(device)}</span>
    </Link>
  );
}
