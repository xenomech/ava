import { Device, DeviceHalo, cn, type DeviceKind } from "@ava/ui";
import {
  TRAIT_BRIGHTNESS,
  TRAIT_COLOR,
  TRAIT_COLOR_TEMP,
  deviceProfile,
  emitsLight,
  isOn,
  numberOf,
  supports,
  traitValue,
  type DeviceDto,
} from "@ava/contracts";
import { PlugZapIcon } from "lucide-react";

import { kelvinToCss } from "@/shared/lib/kelvin";

export function deviceKind(device: DeviceDto): DeviceKind {
  return deviceProfile(device).kind;
}

export function deviceLevel(device: DeviceDto): number {
  if (!isOn(device)) return 0;
  if (!supports(device, TRAIT_BRIGHTNESS)) return 100;

  return numberOf(device, TRAIT_BRIGHTNESS) ?? 100;
}

export function OnAPlug({ className }: { className?: string }) {
  return (
    <span
      title="On a smart plug"
      aria-label="On a smart plug"
      className={cn(
        "grid size-6 place-items-center rounded-full border border-border bg-surface text-subtle",
        className,
      )}
    >
      <PlugZapIcon className="size-3.5" aria-hidden />
    </span>
  );
}

export function deviceColor(device: DeviceDto): string {
  const picked = traitValue(device, TRAIT_COLOR);
  if (typeof picked === "string" && picked !== "") return picked;

  return kelvinToCss(numberOf(device, TRAIT_COLOR_TEMP) ?? 2700);
}

export function DeviceStage({
  devices,
  focusedID,
  levelOverride,
  className,
}: {
  devices: DeviceDto[];
  focusedID: string;
  levelOverride?: number | null;
  className?: string;
}) {
  return (
    <div className={cn("relative grid min-h-0 flex-1 place-items-center", className)}>
      {devices.map((device) => {
        const isFocused = device.id === focusedID;
        const level = isFocused && levelOverride != null ? levelOverride : deviceLevel(device);
        const color = deviceColor(device);

        return (
          <div
            key={device.id}
            aria-hidden={!isFocused}
            className={cn(
              "absolute inset-0 grid place-items-center",
              "transition-[opacity,transform] duration-[380ms] ease-out-soft",
              isFocused ? "opacity-100 scale-100" : "pointer-events-none opacity-0 scale-90",
            )}
            style={{ "--level": level, "--lit": color } as React.CSSProperties}
          >
            {emitsLight(deviceKind(device)) ? <DeviceHalo className="w-[46%]" /> : null}
            {device.kind === "plug" && device.appliance ? (
              <OnAPlug className="absolute right-0 top-0" />
            ) : null}
            <Device
              kind={deviceKind(device)}
              level={level}
              color={color}
              className="h-[58%] max-h-[380px]"
            />
          </div>
        );
      })}
    </div>
  );
}
