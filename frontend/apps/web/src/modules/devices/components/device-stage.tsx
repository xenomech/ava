import { Device, DeviceHalo, cn } from "@ava/ui";
import { useCallback, useState } from "react";

import { useDeviceGesture } from "@/shared/hooks/use-device-gesture";
import { isBehindOverlay } from "@/shared/lib/overlay";
import { useDeviceStore, useDevices, useFocusedDevice } from "../store";

export function DeviceStage({ onHold, className }: { onHold?: () => void; className?: string }) {
  const { devices, focused } = useDevices();
  const { device, setLevel, toggle, step } = useFocusedDevice();

  const [dragLevel, setDragLevel] = useState<number | null>(null);
  const startLevel = useDeviceStore((s) => s.devices.find((d) => d.id === s.focused)?.level ?? 0);

  const gesture = useDeviceGesture({
    onTap: useCallback(() => toggle(device.id), [toggle, device.id]),
    onHold: useCallback(() => onHold?.(), [onHold]),
    onDim: useCallback((delta) => setDragLevel(startLevel + delta), [startLevel]),
    onDimEnd: useCallback(() => {
      if (dragLevel !== null) setLevel(device.id, dragLevel);
      setDragLevel(null);
    }, [dragLevel, setLevel, device.id]),
    onSwipe: useCallback((direction) => step(direction), [step]),
  });

  const shown = dragLevel ?? device.level;

  return (
    <div
      {...gesture}
      // oxlint-disable-next-line prefer-tag-over-role
      role="slider"
      tabIndex={0}
      aria-label={`${device.name} brightness`}
      aria-valuemin={0}
      aria-valuemax={100}
      aria-valuenow={Math.round(shown)}
      onKeyDown={(event) => {
        if (isBehindOverlay(event.currentTarget)) return;

        const stepBy = event.shiftKey ? 10 : 5;
        if (event.key === "ArrowUp") setLevel(device.id, device.level + stepBy);
        else if (event.key === "ArrowDown") setLevel(device.id, device.level - stepBy);
        else if (event.key === "ArrowRight") step(1);
        else if (event.key === "ArrowLeft") step(-1);
        else if (event.key === "Enter" || event.key === " ") toggle(device.id);
        else return;
        event.preventDefault();
      }}
      className={cn(
        "relative grid min-h-0 flex-1 cursor-grab place-items-center outline-none",
        "active:cursor-grabbing",
        className,
      )}
    >
      {devices.map((d) => {
        const isFocused = d.id === focused;
        const level = isFocused ? shown : d.level;

        return (
          <div
            key={d.id}
            aria-hidden={!isFocused}
            className={cn(
              "absolute inset-0 grid place-items-center",
              "transition-[opacity,transform] duration-[380ms] ease-out-soft",
              isFocused ? "opacity-100 scale-100" : "pointer-events-none opacity-0 scale-90",
            )}
            style={{ "--level": level, "--lit": d.color } as React.CSSProperties}
          >
            <DeviceHalo className="w-[46%]" />
            <Device kind={d.kind} level={level} color={d.color} className="h-[58%] max-h-[380px]" />
          </div>
        );
      })}
    </div>
  );
}
