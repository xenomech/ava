import { cn } from "@ava/ui";

import { sceneColor, type ScenePreview } from "../lib/capture";

/** A scene drawn as the light it makes; `live` lights the pane rather than a badge. */
export function SceneLights({
  preview,
  live,
  className,
}: {
  preview: ScenePreview[];
  live: boolean;
  className?: string;
}) {
  const color = sceneColor(preview, "transparent");

  return (
    <span
      className={cn(
        "relative grid h-10 place-items-center overflow-hidden rounded-md bg-bg",
        "transition-[background-image] duration-500 ease-out",
        className,
      )}
      // `color-mix`, not a hex alpha suffix: colours arrive as `rgb()` too.
      style={
        live && color !== "transparent"
          ? {
              backgroundImage: `radial-gradient(78% 110% at 50% 58%, color-mix(in srgb, ${color} 22%, transparent), transparent 72%)`,
            }
          : undefined
      }
    >
      <span className="flex items-center justify-center gap-1.5">
        {preview.map((entry) => (
          <Lamp key={entry.id} entry={entry} />
        ))}
      </span>
    </span>
  );
}

function Lamp({ entry }: { entry: ScenePreview }) {
  if (entry.color === null) {
    return <span aria-hidden className="size-2 rounded-full border border-off/70" />;
  }

  // Floored well above nothing, so a dim lamp still reads as a lamp.
  const size = 6 + Math.round((Math.min(Math.max(entry.level, 0), 100) / 100) * 4);

  return (
    <span
      aria-hidden
      className="rounded-full"
      style={{
        width: size,
        height: size,
        background: entry.color,
        // Tight. A wide glow turns three lamps into one smudge at this size.
        boxShadow: `0 0 5px -1px ${entry.color}`,
      }}
    />
  );
}
