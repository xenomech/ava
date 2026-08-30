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

import { useDevices } from "../hooks/use-devices";

export function CommandPalette() {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();

  const { devices } = useDevices();

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
                void navigate({ to: "/", search: { device: device.id } });
              })}
            >
              {device.name}
              <CommandHint>{device.room || "No room"}</CommandHint>
            </CommandItem>
          ))}
        </CommandGroup>

        <CommandGroup>
          {(
            [
              { label: "Home", to: "/" },
              { label: "Hubs", to: "/settings/hubs" },
              { label: "People", to: "/settings/people" },
              { label: "Appearance", to: "/settings/appearance" },
              { label: "Account", to: "/settings/account" },
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
