import { cn, playSound } from "@ava/ui";
import { TRAIT_POWER, isOn, supports, type DeviceDto, type SceneDto } from "@ava/contracts";
import { PlusIcon } from "lucide-react";
import { useEffect, useLayoutEffect, useMemo, useRef, useState } from "react";

import { matches, scenePreview, type ScenePreview } from "../lib/capture";
import { SceneLights } from "./scene-lights";
import { SceneSheet } from "./scene-sheet";

/** Card width in px. The rail's end padding is derived from it. */
const CARD = 104;

/** What the switch means: the card in the middle is the scene it will play. */
export function SceneRow({
  roomId,
  roomName,
  devices,
  scenes,
  scenesReady,
  armedId,
  onArm,
  onApply,
}: {
  roomId: string;
  roomName: string;
  /** The devices in this room, for previews and for what a new scene captures. */
  devices: DeviceDto[];
  scenes: SceneDto[];
  /** Whether the scene list has arrived, so the rail knows when to take aim. */
  scenesReady: boolean;
  /** `null` is the default: everything on. */
  armedId: string | null;
  /** The centred card changed. Aims the switch; changes nothing in the room. */
  onArm: (scene: SceneDto | null) => void;
  /** Play this scene now. */
  onApply: (scene: SceneDto | null) => void;
}) {
  const [open, setOpen] = useState(false);
  const rail = useRef<HTMLFieldSetElement>(null);

  // Raised while centring on mount, so that scroll cannot arm a neighbouring card.
  const settling = useRef(true);

  // `null` is "All on" and "add" is the save card; only the scenes between arm.
  const stops: (SceneDto | null | "add")[] = [null, ...scenes, "add"];
  const stopId = (stop: (typeof stops)[number]) =>
    stop === "add" ? "add" : stop === null ? "all" : stop.id;

  const armedStop = armedId ?? "all";
  const [centred, setCentred] = useState(armedStop);
  // What the listener last settled on, so a repeat scroll event is a no-op.
  const settledOn = useRef(armedStop);

  // O(scenes × devices), so it is recomputed only when those change, not on scroll.
  const cardFacts = useMemo(() => {
    const switchable = devices.filter((device) => supports(device, TRAIT_POWER));
    const everythingOn = switchable.length > 0 && switchable.every(isOn);

    const map = new Map<string, { preview: ScenePreview[]; live: boolean }>();
    map.set("all", { preview: scenePreview(null, devices), live: everythingOn });

    for (const scene of scenes) {
      map.set(scene.id, { preview: scenePreview(scene, devices), live: matches(scene, devices) });
    }

    return map;
  }, [scenes, devices]);

  // Centre the armed card, keyed on the scene list arriving rather than on arming.
  useLayoutEffect(() => {
    const node = rail.current;
    const card = node?.querySelector<HTMLElement>(`[data-stop="${armedStop}"]`);

    // scrollIntoView, not offsetLeft: that is measured from the room, not the rail.
    if (node && card && centred !== armedStop) {
      settling.current = true;
      card.scrollIntoView({ behavior: "instant", inline: "center", block: "nearest" });
      settledOn.current = armedStop;
      setCentred(armedStop);
    }

    // Cleared on every run, or a run with nothing to do leaves the rail deaf.
    let second = 0;
    const first = requestAnimationFrame(() => {
      second = requestAnimationFrame(() => {
        settling.current = false;
      });
    });

    return () => {
      cancelAnimationFrame(first);
      cancelAnimationFrame(second);
    };
    // `centred` and `armedStop` say where to go, not when; listing them re-enters.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [roomId, scenesReady]);

  // Measured with rects: offsetLeft and scrollLeft have different origins.
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

      // Most scroll events land on the card that is already centred.
      if (settled === settledOn.current) return;

      settledOn.current = settled;
      setCentred(settled);

      // The save card is a stop but not a scene, so the switch keeps its aim.
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
        {/* The bleed is on a wrapper: the rail's percentage padding needs this width. */}
        <div className="-mx-5 w-[calc(100%+2.5rem)] sm:-mx-6 sm:w-[calc(100%+3rem)]">
          {/* A fieldset natively means one choice; `min-w-0` lets it scroll. */}
          <fieldset
            ref={rail}
            className={cn(
              "no-scrollbar flex w-full min-w-0 snap-x snap-mandatory",
              "items-start gap-2.5 overflow-x-auto",
            )}
            // Half the rail less half a card, so the end stops can reach the middle.
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

              const facts = cardFacts.get(id);

              return (
                <SceneCard
                  key={id}
                  id={id}
                  label={stop === null ? "All on" : stop.name}
                  preview={facts?.preview ?? []}
                  centred={isCentred}
                  live={facts?.live ?? false}
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
      // The room changing is worth more than a press click; the delegate gives both.
      data-sound="none"
      onClick={onPick}
      style={{ width: CARD }}
      className={cn(
        "grid shrink-0 snap-center gap-2 rounded-xl border p-2",
        "transition-[transform,opacity,border-color,background-color,box-shadow]",
        "duration-200 ease-out",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
        // Chosen by weight, not by outline: a bright edge shouts what size says.
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
      {/* A dashed pane, so the silhouette matches and only the material says empty. */}
      <span className="grid h-10 place-items-center rounded-md border border-dashed border-border-strong">
        <PlusIcon className="size-4" aria-hidden />
      </span>

      <span className="truncate text-center text-caption font-medium">New scene</span>
    </button>
  );
}
