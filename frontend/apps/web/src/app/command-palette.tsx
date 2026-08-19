import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandHint,
  CommandInput,
  CommandItem,
  CommandList,
} from "@ava/ui";
import { useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";

import { useDeviceStore } from "@/modules/devices";

export function CommandPalette({ onAddDevice }: { onAddDevice: () => void }) {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();

  const devices = useDeviceStore((state) => state.devices);
  const focus = useDeviceStore((state) => state.focus);

  useEffect(() => {
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key.toLowerCase() !== "k" || !(event.metaKey || event.ctrlKey)) return;

      event.preventDefault();
      setOpen((current) => !current);
    };

    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, []);

  const run = (effect: () => void) => () => {
    setOpen(false);
    effect();
  };

  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <CommandInput placeholder="Search devices, rooms, pages…" />

      <CommandList>
        <CommandEmpty>Nothing matches that.</CommandEmpty>

        <CommandGroup>
          {devices.map((device) => (
            <CommandItem
              key={device.id}
              value={`${device.name} ${device.room}`}
              onSelect={run(() => {
                focus(device.id);
                void navigate({ to: "/" });
              })}
            >
              {device.name}
              <CommandHint>{device.room}</CommandHint>
            </CommandItem>
          ))}
        </CommandGroup>

        <CommandGroup>
          <CommandItem value="Add a device" onSelect={run(onAddDevice)}>
            Add a device
            <CommandHint>Pair</CommandHint>
          </CommandItem>

          {(
            [
              { label: "Console", to: "/" },
              { label: "Rooms", to: "/rooms" },
              { label: "Energy", to: "/energy" },
              { label: "People", to: "/settings/members" },
              { label: "Settings", to: "/settings" },
            ] as const
          ).map(({ label, to }) => (
            <CommandItem key={to} value={label} onSelect={run(() => void navigate({ to }))}>
              {label}
              <CommandHint>Go to</CommandHint>
            </CommandItem>
          ))}
        </CommandGroup>
      </CommandList>
    </CommandDialog>
  );
}
