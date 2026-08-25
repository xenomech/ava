import { cn, playSound } from "@ava/ui";
import { TRAIT_POWER, isOn, supports, type DeviceDto, type SceneDto } from "@ava/contracts";
import { PlusIcon } from "lucide-react";
import { useEffect, useLayoutEffect, useRef, useState } from "react";

import { matches, scenePreview, type ScenePreview } from "../capture";
import { SceneLights } from "./scene-lights";
import { SceneSheet } from "./scene-sheet";

/** Card width in px. The rail's end padding is derived from it. */
const CARD = 104;

/**
 * What the switch means, as a carousel it points at.
 *
 * The card in the middle is the one the switch will play. That is the whole
 * model: no separate arming control and nothing to press to choose — you scroll
 * the room to the mood you want and throw the switch. Nothing points at the
 * chosen card, because it is already the only one at full size and full
 * strength; an arrow would only label what the eye has found already.
 *
 * Scrolling arms without applying, which would normally be a dead gesture, so
 * the room previews instead: the glow behind the plate and the paddle both take
 * the centred scene's colour as it passes. The swipe visibly does something, it
 * just does it to the switch rather than to the lights. Tapping the centred
 * card is the shortcut for anyone who would rather not flick at all.
 *
 * A card draws its scene rather than naming it, and the card whose scene is
 * actually running has its pane lit — the same thing the room is doing behind
 * the plate two inches above, so it needs no legend.
 */
export function SceneRow({
  roomId,
  roomName,
  devices,
  scenes,
  armedId,
  onArm,
  onApply,
}: {
  roomId: string;
  roomName: string;
  /** The devices in this room, for previews and for what a new scene captures. */
  devices: DeviceDto[];
  scenes: SceneDto[];
  /** `null` is the default: everything on. */
  armedId: string | null;
  /** The centred card changed. Aims the switch; changes nothing in the room. */
  onArm: (scene: SceneDto | null) => void;
  /** Play this scene now. */
  onApply: (scene: SceneDto | null) => void;
}) {
  const [open, setOpen] = useState(false);
  const rail = useRef<HTMLFieldSetElement>(null);

  /* Raised while the rail is being put where it belongs on mount.
   *
   * Centring the armed card is itself a scroll, and the scroll listener is what
   * decides which card is armed — so on a frame where the cards have not
   * finished sizing, that first programmatic scroll could settle on a
   * neighbour, arm it, and overwrite the choice it was in the middle of
   * restoring. The room came back pointed at the wrong scene, occasionally, for
   * no reason anyone could see. */
  const settling = useRef(true);

  /* `null` is "All on" and "add" is the trailing save card. Neither is a scene
     id, and only the scenes in between can be armed. */
  const stops: (SceneDto | null | "add")[] = [null, ...scenes, "add"];
  const stopId = (stop: (typeof stops)[number]) =>
    stop === "add" ? "add" : stop === null ? "all" : stop.id;

  const armedStop = armedId ?? "all";
  const [centred, setCentred] = useState(armedStop);

  const switchable = devices.filter((device) => supports(device, TRAIT_POWER));
  const everythingOn = switchable.length > 0 && switchable.every(isOn);

  /* Put the armed card under the caret before the first paint, so the rail does
     not visibly slide into place every time the room is opened. */
  useLayoutEffect(() => {
    const node = rail.current;
    const card = node?.querySelector<HTMLElement>(`[data-stop="${armedStop}"]`);
    if (!node || !card) return;

    /* scrollIntoView rather than arithmetic on offsetLeft, which is measured
       from the nearest positioned ancestor — the room, not the rail — and left
       the armed card sitting a little to one side of the caret. */
    settling.current = true;
    card.scrollIntoView({ behavior: "instant", inline: "center", block: "nearest" });
    setCentred(armedStop);

    /* Two frames: one for the scroll to be applied, one for its event to have
       been delivered. */
    requestAnimationFrame(() =>
      requestAnimationFrame(() => {
        settling.current = false;
      }),
    );
    // Only when the room changes. Scrolling drives it from then on.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [roomId]);

  /* Measured with rects, not offsetLeft. offsetLeft is relative to the nearest
     positioned ancestor and scrollLeft to the rail's own content box, and
     comparing the two pins the answer to the first card forever. */
  useEffect(() => {
    const node = rail.current;
    if (!node) return;

    const settle = () => {
      if (settling.current) return;

      const box = node.getBoundingClientRect();
      const middle = box.left + box.width / 2;

      let nearestId: string | null = null;
      let nearestGap = Infinity;

      node.querySelectorAll<HTMLElement>("[data-stop]").forEach((card) => {
        const rect = card.getBoundingClientRect();
        const gap = Math.abs(rect.left + rect.width / 2 - middle);

        if (gap < nearestGap) {
          nearestGap = gap;
          nearestId = card.dataset.stop ?? null;
        }
      });

      if (nearestId === null) return;

      const settled: string = nearestId;
      setCentred(settled);

      /* The save card is a stop on the rail but not a scene, so passing over it
         leaves the switch aimed wherever it already was. */
      if (settled === "add") return;

      onArm(settled === "all" ? null : (scenes.find((scene) => scene.id === settled) ?? null));
    };

    node.addEventListener("scroll", settle, { passive: true });

    return () => node.removeEventListener("scroll", settle);
  }, [scenes, onArm]);

  const centre = (id: string) => {
    rail.current
      ?.querySelector<HTMLElement>(`[data-stop="${id}"]`)
      ?.scrollIntoView({ behavior: "smooth", inline: "center", block: "nearest" });
  };

  return (
    <>
      <div className="grid w-full grid-cols-[minmax(0,1fr)] justify-items-center">
        {/* The bleed lives on a wrapper rather than on the rail itself. The
            rail's end padding is a percentage, and a percentage padding resolves
            against the containing block — so while the rail was the element
            carrying `-mx-5 w-[100%+2.5rem]`, its own 50% was half of the
            narrower parent and every card settled exactly one margin off centre,
            under nothing. */}
        <div className="-mx-5 w-[calc(100%+2.5rem)] sm:-mx-6 sm:w-[calc(100%+3rem)]">
          {/* A fieldset because this natively means "these controls are one
              choice", and the choice has a name worth reading out. `min-w-0`
              undoes its default min-inline-size, which would otherwise refuse
              to let the rail scroll and widen the page instead. */}
          <fieldset
            ref={rail}
            className={cn(
              "no-scrollbar flex w-full min-w-0 snap-x snap-mandatory",
              "items-start gap-2.5 overflow-x-auto",
            )}
            /* Half the rail less half a card, so the first and last stops can
               reach the middle. A fixed inset leaves the ends unable to centre
               and the caret pointing at nothing. */
            style={{ paddingInline: `calc(50% - ${CARD / 2}px)` }}
          >
            <legend className="sr-only">What the {roomName} switch turns on</legend>

            {stops.map((stop) => {
              const id = stopId(stop);
              const isCentred = id === centred;

              if (stop === "add") {
                return (
                  <SaveCard
                    key="add"
                    centred={isCentred}
                    roomName={roomName}
                    onOpen={() => (isCentred ? setOpen(true) : centre("add"))}
                  />
                );
              }

              return (
                <SceneCard
                  key={id}
                  id={id}
                  label={stop === null ? "All on" : stop.name}
                  preview={scenePreview(stop, devices)}
                  centred={isCentred}
                  live={stop === null ? everythingOn : matches(stop, devices)}
                  onPick={() => {
                    if (!isCentred) {
                      centre(id);

                      return;
                    }

                    playSound("on");
                    onApply(stop);
                  }}
                />
              );
            })}
          </fieldset>
        </div>
      </div>

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

function SceneCard({
  id,
  label,
  preview,
  centred,
  live,
  onPick,
}: {
  id: string;
  label: string;
  preview: ScenePreview[];
  centred: boolean;
  live: boolean;
  onPick: () => void;
}) {
  return (
    <button
      type="button"
      data-stop={id}
      aria-pressed={centred}
      /* The room changing is worth more than a press click, and the delegated
         listener would otherwise give it both. */
      data-sound="none"
      onClick={onPick}
      style={{ width: CARD }}
      className={cn(
        "grid shrink-0 snap-center gap-2 rounded-xl border p-2",
        "transition-[transform,opacity,border-color,background-color,box-shadow]",
        "duration-200 ease-out",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
        /* Chosen by weight, not by outline. The old white border was the one
           bright edge on a screen whose every other border is a dark hairline,
           and it shouted across the room to say something the size and the
           brightness of the card were already saying. */
        centred
          ? cn(
              "scale-100 border-border-strong bg-raised opacity-100",
              "shadow-[0_10px_24px_-14px_rgb(0_0_0/0.9)]",
            )
          : "scale-[0.88] border-border bg-surface opacity-40 shadow-none",
      )}
    >
      <SceneLights preview={preview} live={live} />

      <span className="min-w-0 truncate text-center text-caption font-medium">
        {label}
        {live ? <span className="sr-only"> (showing now)</span> : null}
      </span>
    </button>
  );
}

/** The end of the rail: the same card, waiting to be filled in. */
function SaveCard({
  centred,
  roomName,
  onOpen,
}: {
  centred: boolean;
  roomName: string;
  onOpen: () => void;
}) {
  return (
    <button
      type="button"
      data-stop="add"
      onClick={onOpen}
      aria-label={`Scenes in ${roomName}`}
      style={{ width: CARD }}
      className={cn(
        "grid shrink-0 snap-center gap-2 rounded-xl border p-2",
        "transition-[transform,opacity,border-color,background-color] duration-200 ease-out",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
        centred
          ? "scale-100 border-border-strong bg-raised text-fg opacity-100"
          : "scale-[0.88] border-border bg-surface text-muted opacity-40",
      )}
    >
      {/* A dashed pane where the lights would be, so the silhouette matches its
          neighbours and only the material says it is empty. */}
      <span className="grid h-10 place-items-center rounded-md border border-dashed border-border-strong">
        <PlusIcon className="size-4" aria-hidden />
      </span>

      <span className="truncate text-center text-caption font-medium">New scene</span>
    </button>
  );
}
