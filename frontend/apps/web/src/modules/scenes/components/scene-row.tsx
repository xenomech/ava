import { cn, playSound } from "@ava/ui";
import { TRAIT_POWER, isOn, supports, type DeviceDto, type SceneDto } from "@ava/contracts";
import { PlusIcon } from "lucide-react";
import { useState } from "react";

import { matches } from "../capture";
import { SceneSheet } from "./scene-sheet";

/**
 * What the switch means, as a row of chips under it.
 *
 * A scene here is not a shortcut sitting beside the switch — it is the switch's
 * setting. Flicking up plays whichever chip is filled, so the row answers the
 * only question a switch can raise once a room has more than one mood: on, but
 * on to what?
 *
 * Tapping a chip applies it there and then rather than merely arming it. A tap
 * that changes nothing visible reads as a dead control, and "switch the room to
 * Evening" is what someone tapping "Evening" wanted in the first place; that it
 * also aims the switch is the part they get for free.
 *
 * The lit dot marks the scene the room is currently matching, which is a
 * different fact from the one the chip fill carries. You can be aimed at
 * Evening while the room is plainly not in it — you turned a lamp off by hand —
 * and the row should say so rather than quietly claim otherwise.
 */
export function SceneRow({
  roomId,
  roomName,
  devices,
  scenes,
  armedId,
  onPick,
}: {
  roomId: string;
  roomName: string;
  /** The devices in this room, for matching and for what a new scene captures. */
  devices: DeviceDto[];
  scenes: SceneDto[];
  /** `null` is the default: everything on. */
  armedId: string | null;
  onPick: (scene: SceneDto | null) => void;
}) {
  const [open, setOpen] = useState(false);

  /* "All on" is live when everything the switch would actually reach is on.
     Measured over the devices that have a power trait, not all of them, so a
     sensor sitting in the room cannot make the state unreachable. */
  const switchable = devices.filter((device) => supports(device, TRAIT_POWER));
  const everythingOn = switchable.length > 0 && switchable.every(isOn);

  return (
    <>
      {/* A fieldset because that is the native "these controls are one choice"
          element, and the choice does have a name worth reading out. `min-w-0`
          undoes its default min-inline-size, which would otherwise refuse to
          let the row scroll and widen the page instead. */}
      <fieldset
        className={cn(
          "no-scrollbar -mx-5 flex w-[calc(100%+2.5rem)] min-w-0 items-center gap-2",
          "overflow-x-auto px-5 sm:-mx-6 sm:w-[calc(100%+3rem)] sm:px-6",
          /* `safe` centring, not plain `center`. A centred flex row that
             overflows pushes its first item off the left edge where nothing can
             scroll back to it; `safe` falls back to start in exactly that
             case. */
          "[justify-content:safe_center]",
        )}
      >
        <legend className="sr-only">What the {roomName} switch turns on</legend>

        <SceneChip
          label="All on"
          armed={armedId === null}
          live={everythingOn}
          onPick={() => onPick(null)}
        />

        {scenes.map((scene) => (
          <SceneChip
            key={scene.id}
            label={scene.name}
            armed={scene.id === armedId}
            live={matches(scene, devices)}
            onPick={() => onPick(scene)}
          />
        ))}

        <button
          type="button"
          onClick={() => setOpen(true)}
          aria-label={`Scenes in ${roomName}`}
          className={cn(
            "grid size-8 shrink-0 place-items-center rounded-full border border-dashed",
            "border-border-strong text-muted transition-colors duration-150 ease-out",
            "hover:border-fg hover:text-fg focus-visible:outline-none focus-visible:ring-2",
            "focus-visible:ring-fg",
          )}
        >
          <PlusIcon className="size-3.5" aria-hidden />
        </button>
      </fieldset>

      <SceneSheet
        open={open}
        onOpenChange={setOpen}
        roomId={roomId}
        roomName={roomName}
        devices={devices}
        scenes={scenes}
      />
    </>
  );
}

function SceneChip({
  label,
  armed,
  live,
  onPick,
}: {
  label: string;
  armed: boolean;
  live: boolean;
  onPick: () => void;
}) {
  return (
    <button
      type="button"
      aria-pressed={armed}
      /* The room changing is worth more than a press click, and the delegated
         listener would otherwise give it both. */
      data-sound="none"
      onClick={() => {
        playSound("on");
        onPick();
      }}
      className={cn(
        "flex h-8 shrink-0 items-center gap-1.5 rounded-full border px-3.5",
        "text-caption font-medium transition-colors duration-150 ease-out",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
        "focus-visible:ring-offset-2 focus-visible:ring-offset-bg",
        armed
          ? "border-fg bg-fg text-bg"
          : "border-border bg-surface text-muted hover:border-border-strong hover:text-fg",
      )}
    >
      <span
        aria-hidden
        className={cn(
          "size-1.5 rounded-full transition-colors duration-300 ease-out",
          live ? "bg-success" : "bg-transparent",
        )}
      />
      {label}
      {live ? <span className="sr-only">(showing now)</span> : null}
    </button>
  );
}
