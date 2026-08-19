import { Button, Device, Slider, Switch, cn } from "@ava/ui";
import { useEffect, useRef } from "react";

import { isBehindOverlay } from "@/shared/lib/overlay";
import { ColorControl } from "../components/color-control";
import { DeviceStage } from "../components/device-stage";
import { useDevices, useFocusedDevice } from "../store";

export function ConsolePage() {
  const { devices, focused, focus } = useDevices();
  const { device, setLevel, toggle, setColor, setWhite, step } = useFocusedDevice();

  const root = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      const target = event.target as HTMLElement;
      if (target.matches("input, textarea, [contenteditable]")) return;

      if (isBehindOverlay(root.current)) return;

      const stepBy = event.shiftKey ? 10 : 5;
      if (event.key === "ArrowRight") step(1);
      else if (event.key === "ArrowLeft") step(-1);
      else if (event.key === "ArrowUp") setLevel(device.id, device.level + stepBy);
      else if (event.key === "ArrowDown") setLevel(device.id, device.level - stepBy);
      else if (event.key === " ") toggle(device.id);
      else return;

      event.preventDefault();
    };

    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [device.id, device.level, setLevel, toggle, step]);

  const watts = ((device.watts * device.level) / 100).toFixed(1);

  return (
    <div ref={root} className="grid h-full grid-rows-[minmax(0,1fr)_auto]">
      <div className="min-h-0 overflow-y-auto lg:grid lg:grid-cols-[minmax(0,1fr)_330px] lg:overflow-hidden">
        <main className="grid min-h-0 grid-rows-[34vh_auto] p-5 sm:p-6 lg:grid-rows-[minmax(0,1fr)_auto]">
          <DeviceStage />

          <div className="flex items-end justify-between gap-4 pt-4">
            <div>
              <h1 className="text-display font-semibold">{device.name}</h1>
              <p className="mt-1.5 text-small text-muted">
                {device.room} · <span className="font-mono tabular">{watts}W</span>
              </p>
            </div>

            <div className="flex items-start gap-0.5" data-slot="reading">
              <b className="text-hero font-semibold">{device.level}</b>
              <span className="mt-1 text-small text-subtle">%</span>
            </div>
          </div>
        </main>

        <aside className="flex flex-col gap-5 border-t border-border p-5 lg:min-h-0 lg:overflow-y-auto lg:border-t-0 lg:border-l">
          <div className="grid gap-2.5">
            <div className="flex items-center justify-between">
              <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
                Power
              </span>
              <Switch
                checked={device.level > 0}
                onCheckedChange={() => toggle(device.id)}
                aria-label={`${device.name} power`}
              />
            </div>
          </div>

          <div className="grid gap-2.5">
            <div className="flex items-baseline justify-between">
              <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
                Brightness
              </span>
              <output className="font-mono text-small tabular">{device.level}%</output>
            </div>
            <Slider
              value={[device.level]}
              max={100}
              step={1}
              lit
              aria-label="Brightness"
              onValueChange={([v]) => setLevel(device.id, v ?? 0)}
              style={{ "--lit": device.color } as React.CSSProperties}
            />
          </div>

          <div className="grid gap-2.5">
            <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
              Colour
            </span>
            <ColorControl
              color={device.color}
              kelvin={device.kelvin}
              onWhite={(k) => setWhite(device.id, k)}
              onColor={(c) => setColor(device.id, c)}
            />
          </div>

          <div className="hidden gap-2.5 lg:grid">
            <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
              Shortcuts
            </span>
            <dl className="grid gap-px overflow-hidden rounded-md border border-border bg-border">
              {[
                ["Change device", "← →"],
                ["Brightness", "↑ ↓"],
                ["Toggle", "Space"],
              ].map(([label, keys]) => (
                <div
                  key={label}
                  className="flex items-center justify-between bg-surface px-3.5 py-2.5"
                >
                  <dt className="text-small">{label}</dt>
                  <dd className="font-mono text-caption text-subtle">{keys}</dd>
                </div>
              ))}
            </dl>
          </div>
        </aside>
      </div>

      <footer
        className="flex gap-2 overflow-x-auto border-t border-border p-3"
        aria-label="All devices"
      >
        {devices.map((d) => (
          <Button
            key={d.id}
            variant="ghost"
            aria-current={d.id === focused}
            onClick={() => focus(d.id)}
            className={cn(
              "h-auto w-32 shrink-0 flex-col items-start gap-1.5 rounded-md border-border p-2.5",
              d.id === focused && "border-fg",
            )}
            style={{ "--level": d.level, "--lit": d.color } as React.CSSProperties}
          >
            <span className="grid h-14 w-full place-items-center">
              <Device kind={d.kind} level={d.level} color={d.color} className="h-full" />
            </span>
            <span className="w-full truncate text-left text-small font-semibold">{d.name}</span>
            <span className="font-mono text-caption font-normal text-subtle tabular">
              {d.level > 0 ? `${d.level}%` : "Off"}
            </span>
          </Button>
        ))}
      </footer>
    </div>
  );
}
