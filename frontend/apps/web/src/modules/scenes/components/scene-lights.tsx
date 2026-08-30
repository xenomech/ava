import { cn } from "@ava/ui";

import { sceneColor, type ScenePreview } from "../lib/capture";

/**
 * A scene drawn as the light it makes: one lamp per fixture, sized by
 * brightness, coloured by what that fixture will be, a bare ring for whatever
 * the scene leaves off.
 *
 * `live` means the room is doing this right now, and it is shown by lighting
 * the pane rather than by a badge. A scene that is running has its light on, so
 * the panel behind the lamps glows the way the room does behind the switch — no
 * legend, no coloured dot standing in for a sentence, and the same language the
 * plate already speaks two inches above it.
 */
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
      /* `color-mix` rather than a hex alpha suffix. A scene's colour arrives as
         `rgb(255 216 168)` from the warmth ramp and as `#rrggbb` from a colour
         picker, and appending `2e` to the first produced invalid CSS — which
         browsers drop silently, so the pane simply never lit and nothing said
         why. */
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

  /* Floored well above nothing, so a dim lamp still reads as a lamp rather than
     as dust on the screen. */
  const size = 6 + Math.round((Math.min(Math.max(entry.level, 0), 100) / 100) * 4);

  return (
    <span
      aria-hidden
      className="rounded-full"
      style={{
        width: size,
        height: size,
        background: entry.color,
        /* Tight. A wide glow turns three lamps into one smudge at this size. */
        boxShadow: `0 0 5px -1px ${entry.color}`,
      }}
    />
  );
}
