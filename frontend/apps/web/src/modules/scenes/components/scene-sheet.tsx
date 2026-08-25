import {
  Button,
  Drawer,
  DrawerBody,
  DrawerContent,
  DrawerDescription,
  DrawerTitle,
  Input,
  cn,
} from "@ava/ui";
import type { DeviceDto, SceneDto } from "@ava/contracts";
import { Trash2Icon } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

import { capture, describe } from "../capture";
import { useSceneActions } from "../use-scenes";

/**
 * Saving the room as a scene, and clearing out the ones you no longer want.
 *
 * The saved scene is described in full before it is saved. A scene is invisible
 * state — a name that stands for eleven trait values — and the one moment it
 * can honestly be shown is now, while the room it came from is still on screen
 * behind the sheet.
 *
 * There is no renaming and no reordering. Both are real, and neither is worth a
 * screen yet: a scene that has the wrong name has been saved for about ten
 * seconds, and deleting and saving again costs one more tap than renaming
 * would.
 */
export function SceneSheet({
  open,
  onOpenChange,
  roomId,
  roomName,
  devices,
  scenes,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  roomId: string;
  roomName: string;
  devices: DeviceDto[];
  scenes: SceneDto[];
}) {
  const [name, setName] = useState("");
  const { save, remove } = useSceneActions(roomId);

  /* Recomputed as the room reports, so the list below is the scene that will
     actually be written. Freezing it at open would be steadier to read and
     occasionally a lie. */
  const targets = useMemo(() => capture(devices), [devices]);

  useEffect(() => {
    if (open) setName("");
  }, [open]);

  const named = name.trim();
  const taken = scenes.some((scene) => scene.name.toLowerCase() === named.toLowerCase());
  const saveable = named.length > 0 && !taken && targets.length > 0;

  const submit = () => {
    if (!saveable) return;

    save.mutate({ name: named, targets }, { onSuccess: () => onOpenChange(false) });
  };

  return (
    <Drawer open={open} onOpenChange={onOpenChange}>
      <DrawerContent className="max-h-[85dvh]">
        <DrawerTitle className="px-5 pt-2 text-title font-semibold">Scenes</DrawerTitle>
        <DrawerDescription className="px-5 pb-1 text-small text-muted">
          A scene remembers how {roomName} is set right now. Flicking the switch up plays whichever
          one you have chosen.
        </DrawerDescription>

        <DrawerBody className="grid content-start gap-6 pb-8 pt-3">
          {scenes.length > 0 ? (
            <section className="grid gap-2">
              <Heading>Saved</Heading>

              {scenes.map((scene) => (
                <div
                  key={scene.id}
                  className={cn(
                    "flex min-h-12 items-center gap-3 rounded-lg border border-border",
                    "bg-surface px-3.5 py-2",
                  )}
                >
                  <span className="min-w-0 flex-1 truncate text-small font-semibold">
                    {scene.name}
                  </span>
                  <span className="shrink-0 font-mono text-caption text-subtle tabular">
                    {countOn(scene)} on
                  </span>
                  <button
                    type="button"
                    disabled={remove.isPending && remove.variables === scene.id}
                    onClick={() => remove.mutate(scene.id)}
                    aria-label={`Delete ${scene.name}`}
                    className={cn(
                      "tap grid size-8 shrink-0 place-items-center rounded-md text-subtle",
                      "transition-colors duration-150 ease-out hover:text-danger",
                      "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
                      "disabled:opacity-40",
                    )}
                  >
                    <Trash2Icon className="size-4" aria-hidden />
                  </button>
                </div>
              ))}
            </section>
          ) : null}

          <section className="grid gap-3">
            <Heading>Save this arrangement</Heading>

            <Input
              value={name}
              maxLength={80}
              placeholder="Evening"
              aria-label="Scene name"
              enterKeyHint="done"
              onChange={(event) => setName(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === "Enter") submit();
              }}
            />

            {taken ? (
              <p className="text-caption text-danger">
                {roomName} already has a scene called that.
              </p>
            ) : null}

            <ul className="grid gap-1.5 rounded-md border border-border bg-surface p-3.5">
              {devices.map((device) => (
                <li key={device.id} className="flex items-baseline justify-between gap-4">
                  <span className="min-w-0 truncate text-caption text-subtle">{device.name}</span>
                  <span className="shrink-0 font-mono text-caption text-muted tabular">
                    {describe(device, targets)}
                  </span>
                </li>
              ))}
            </ul>

            <Button onClick={submit} disabled={!saveable || save.isPending}>
              {save.isPending ? "Saving" : "Save scene"}
            </Button>
          </section>
        </DrawerBody>
      </DrawerContent>
    </Drawer>
  );
}

function Heading({ children }: { children: string }) {
  return (
    <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
      {children}
    </span>
  );
}

function countOn(scene: SceneDto): number {
  return scene.targets.filter((target) => target.trait === "power" && target.value === true).length;
}
