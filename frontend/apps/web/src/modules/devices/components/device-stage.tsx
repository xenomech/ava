import { Device, DeviceHalo, cn, type DeviceKind } from "@ava/ui";
import { deviceProfile, type DeviceDto } from "@ava/contracts";

import { kelvinToCss } from "@/shared/lib/kelvin";

export function deviceKind(device: DeviceDto): DeviceKind {
  return deviceProfile(device).kind;
}

export function deviceLevel(device: DeviceDto): number {
  if (!device.state.power) return 0;

  return device.state.brightness ?? 100;
}

export function deviceColor(device: DeviceDto): string {
  return kelvinToCss(device.state.color_temp ?? 2700);
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
            <DeviceHalo className="w-[46%]" />
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
