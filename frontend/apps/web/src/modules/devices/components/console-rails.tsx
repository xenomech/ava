import { Device, cn } from "@ava/ui";
import { isOn, supports, TRAIT_BRIGHTNESS, type DeviceDto } from "@ava/contracts";
import type { CSSProperties } from "react";

import type { RoomGroup } from "../use-room-groups";
import { deviceColor, deviceKind, deviceLevel } from "../device-view";
import { OnAPlug } from "./on-a-plug";

// A room's dot is lit in the colour temperature of whatever is on inside it, so
// the rail reads as a plan of the home rather than a list of words.
function roomLight(group: RoomGroup): { color: string; level: number } {
  const lit = group.devices.filter(isOn);
  if (lit.length === 0) return { color: "var(--color-off)", level: 0 };

  const level = Math.round(lit.reduce((sum, d) => sum + deviceLevel(d), 0) / lit.length);

  return { color: deviceColor(lit[0]!), level };
}

export function RoomRail({
  groups,
  activeKey,
  onPick,
  className,
}: {
  groups: RoomGroup[];
  activeKey: string;
  onPick: (key: string) => void;
  className?: string;
}) {
  return (
    <nav aria-label="Rooms" className={cn("flex flex-col gap-1", className)}>
      {groups.map((group) => {
        const light = roomLight(group);

        return (
          <button
            key={group.key}
            type="button"
            aria-current={group.key === activeKey}
            onClick={() => onPick(group.key)}
            style={{ "--lit": light.color, "--level": light.level } as CSSProperties}
            className={cn(
              "flex items-center gap-3 rounded-sm px-3 py-2.5 text-left",
              "transition-colors duration-150 ease-out",
              "text-muted hover:bg-surface hover:text-fg",
              "aria-[current=true]:bg-raised aria-[current=true]:text-fg",
            )}
          >
            <Dot />
            <span className="min-w-0 flex-1 truncate text-small font-medium">{group.name}</span>
            <span className="text-caption text-subtle tabular">
              {group.on > 0 ? group.on : "—"}
            </span>
          </button>
        );
      })}
    </nav>
  );
}

export function RoomTabs({
  groups,
  activeKey,
  onPick,
  className,
}: {
  groups: RoomGroup[];
  activeKey: string;
  onPick: (key: string) => void;
  className?: string;
}) {
  return (
    <nav
      aria-label="Rooms"
      className={cn("flex gap-1.5 overflow-x-auto [scrollbar-width:none]", className)}
    >
      {groups.map((group) => {
        const light = roomLight(group);

        return (
          <button
            key={group.key}
            type="button"
            aria-current={group.key === activeKey}
            onClick={() => onPick(group.key)}
            style={{ "--lit": light.color, "--level": light.level } as CSSProperties}
            className={cn(
              "flex shrink-0 items-center gap-2 rounded-full border border-border px-3 py-2",
              "text-small font-medium text-muted transition-colors duration-150 ease-out",
              "aria-[current=true]:border-fg aria-[current=true]:text-fg",
            )}
          >
            <Dot />
            {group.name}
          </button>
        );
      })}
    </nav>
  );
}

// The inset hairline keeps a cool-white room readable on the light theme, where
// a 5000K dot is all but the same colour as the page.
function Dot() {
  return (
    <span
      aria-hidden
      className="size-2 shrink-0 rounded-full bg-[var(--lit)] transition-shadow duration-300 ease-out"
      style={{
        boxShadow:
          "0 0 calc(var(--level) / 100 * 12px) var(--lit), inset 0 0 0 1px var(--palette-glass-edge)",
      }}
    />
  );
}

export function FixtureStrip({
  devices,
  focusedID,
  onFocus,
  layout = "row",
  className,
}: {
  devices: DeviceDto[];
  focusedID: string;
  onFocus: (id: string) => void;
  /** "row" scrolls sideways under the hero; "grid" wraps inside the docked panel. */
  layout?: "row" | "grid";
  className?: string;
}) {
  // One fixture in a room needs no picker — the hero is already it.
  if (devices.length < 2) return null;

  return (
    <div
      role="tablist"
      aria-label="Fixtures in this room"
      className={cn(
        layout === "grid"
          ? "grid grid-cols-2 gap-2"
          : "flex gap-2 overflow-x-auto [scrollbar-width:none]",
        className,
      )}
    >
      {devices.map((device) => (
        <Fixture
          key={device.id}
          device={device}
          focused={device.id === focusedID}
          fluid={layout === "grid"}
          onFocus={() => onFocus(device.id)}
        />
      ))}
    </div>
  );
}

function Fixture({
  device,
  focused,
  fluid,
  onFocus,
}: {
  device: DeviceDto;
  focused: boolean;
  fluid: boolean;
  onFocus: () => void;
}) {
  const level = deviceLevel(device);

  return (
    <button
      type="button"
      role="tab"
      aria-selected={focused}
      onClick={onFocus}
      style={{ "--level": level, "--lit": deviceColor(device) } as CSSProperties}
      className={cn(
        "relative grid gap-1 rounded-lg border border-border p-2.5 text-left",
        fluid ? "w-full min-w-0" : "w-[104px] shrink-0",
        "transition-colors duration-150 ease-out hover:border-border-strong",
        "aria-selected:border-fg aria-selected:bg-surface",
        device.status === "offline" && "opacity-50",
      )}
    >
      {device.kind === "plug" && device.appliance ? (
        <OnAPlug className="absolute right-1.5 top-1.5 size-5" />
      ) : null}

      <span className="grid h-12 place-items-center">
        <Device
          kind={deviceKind(device)}
          level={level}
          color={deviceColor(device)}
          className="h-full"
        />
      </span>

      <span className="truncate text-caption font-semibold">{device.name}</span>
      <span className="text-micro text-muted tabular">{status(device, level)}</span>
    </button>
  );
}

function status(device: DeviceDto, level: number) {
  if (device.status === "offline") return "Offline";
  if (!isOn(device)) return "Off";
  if (!supports(device, TRAIT_BRIGHTNESS)) return "On";

  return `${level}%`;
}
