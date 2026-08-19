import { Chip, cn } from "@ava/ui";

import { useDevices } from "../store";

const HOURS = [
  38, 52, 30, 64, 48, 72, 90, 56, 44, 61, 35, 50, 44, 68, 78, 52, 40, 58, 84, 66, 48, 36, 30, 26,
];

export function EnergyPage() {
  const { devices } = useDevices();
  const peak = Math.max(...HOURS);

  const draw = devices.reduce((sum, d) => sum + (d.watts * d.level) / 100, 0);

  const byDevice = [...devices]
    .map((d) => ({ ...d, kwh: (d.watts * d.level) / 100 / 100 + d.watts / 12 }))
    .sort((a, b) => b.kwh - a.kwh);

  const most = byDevice[0]?.kwh || 1;

  return (
    <div className="grid gap-6 p-5 sm:p-6">
      <header className="grid gap-4 sm:flex sm:items-end sm:justify-between">
        <div>
          <h1 className="text-display font-semibold">Energy</h1>
          <p className="mt-1 text-small text-muted">Whole home · last 24 hours</p>
        </div>
        <div className="flex gap-2">
          {["Day", "Week", "Month"].map((range, i) => (
            <Chip key={range} tone={i === 0 ? "neutral" : "muted"}>
              {range}
            </Chip>
          ))}
        </div>
      </header>

      <div className="grid grid-cols-[repeat(auto-fit,minmax(150px,1fr))] gap-3">
        {[
          ["Today", "6.2 kWh"],
          ["Cost", "₹52"],
          ["Current draw", `${draw.toFixed(1)} W`],
          ["vs yesterday", "−12%"],
        ].map(([label, value]) => (
          <div key={label} className="rounded-md border border-border bg-surface p-4">
            <div className="text-micro font-semibold uppercase tracking-caps text-subtle">
              {label}
            </div>
            <div className="mt-1.5 text-display font-semibold tabular">{value}</div>
          </div>
        ))}
      </div>

      <figure className="grid gap-3 rounded-lg border border-border bg-surface p-5">
        <figcaption className="sr-only">Consumption by hour</figcaption>

        <div className="flex h-56 items-end gap-1.5">
          {HOURS.map((height, i) => (
            <span
              key={i}
              className={cn("flex-1 rounded-t-[3px]", height === peak ? "bg-fg" : "bg-subtle")}
              style={{ height: `${height}%` }}
            />
          ))}
        </div>

        <div className="flex justify-between font-mono text-caption text-subtle">
          {["00:00", "06:00", "12:00", "18:00", "23:59"].map((t) => (
            <span key={t}>{t}</span>
          ))}
        </div>
      </figure>

      <section className="overflow-hidden rounded-lg border border-border">
        <h2 className="border-b border-border bg-surface px-4 py-3 text-caption font-semibold uppercase tracking-caps text-subtle">
          By device
        </h2>

        <ul>
          {byDevice.map((device) => (
            <li
              key={device.id}
              className="flex items-center gap-4 border-b border-border bg-surface px-4 py-3 last:border-b-0"
            >
              <span className="w-40 shrink-0">
                <span className="block text-body font-medium">{device.name}</span>
                <span className="block text-caption text-subtle">{device.room}</span>
              </span>

              <span className="h-1 flex-1 overflow-hidden rounded-[2px] bg-raised">
                <span
                  className="block h-full bg-fg"
                  style={{ width: `${(device.kwh / most) * 100}%` }}
                />
              </span>

              <span className="w-20 shrink-0 text-right font-mono text-small tabular">
                {device.kwh.toFixed(1)} kWh
              </span>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}
