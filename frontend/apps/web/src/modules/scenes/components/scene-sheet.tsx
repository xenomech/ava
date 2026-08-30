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
import { useMemo, useState } from "react";

import { capture, describe, matches, scenePreview } from "../lib/capture";
import { SceneLights } from "./scene-lights";
import { useSceneActions } from "../hooks/use-scenes";

/** The scenes a room has, and the one it is about to get. */
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

  // Recomputed as the room reports, so the panel is the scene that will be written.
  const targets = useMemo(() => capture(devices), [devices]);
  const preview = useMemo(() => scenePreview(null, devices), [devices]);

  // Cleared on close rather than on open, so no effect and no frame of the old name.
  const openChange = (next: boolean) => {
    if (!next) setName("");
    onOpenChange(next);
  };

  const named = name.trim();
  const taken = scenes.some((scene) => scene.name.toLowerCase() === named.toLowerCase());
  const saveable = named.length > 0 && !taken && targets.length > 0;

  const submit = () => {
    if (!saveable) return;

    save.mutate({ name: named, targets }, { onSuccess: () => openChange(false) });
  };

  return (
    <Drawer open={open} onOpenChange={openChange}>
      <DrawerContent className="max-h-[85dvh]">
        <DrawerTitle className="px-5 pt-2 text-title font-semibold">Scenes</DrawerTitle>
        <DrawerDescription className="px-5 pb-1 text-small text-muted">
          A scene remembers how {roomName} is set. Flicking the switch up plays whichever one you
          have scrolled to.
        </DrawerDescription>

        <DrawerBody className="grid content-start gap-5 pb-8 pt-4">
          {scenes.length > 0 ? (
            // Hairlines rather than a border each: a divided list reads as one thing.
            <ul className="grid divide-y divide-border border-y border-border">
              {scenes.map((scene) => (
                <li key={scene.id} className="flex items-center gap-3 py-2.5">
                  <SceneLights
                    preview={scenePreview(scene, devices)}
                    live={matches(scene, devices)}
                    className="w-[72px] shrink-0"
                  />

                  <span className="min-w-0 flex-1 truncate text-small font-medium">
                    {scene.name}
                  </span>

                  <button
                    type="button"
                    disabled={remove.isPending && remove.variables === scene.id}
                    onClick={() => remove.mutate(scene.id)}
                    aria-label={`Delete ${scene.name}`}
                    className={cn(
                      "tap grid size-9 shrink-0 place-items-center rounded-md text-subtle",
                      "transition-colors duration-150 ease-out hover:text-danger",
                      "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
                      "disabled:opacity-40",
                    )}
                  >
                    <Trash2Icon className="size-4" aria-hidden />
                  </button>
                </li>
              ))}
            </ul>
          ) : null}

          <section className="grid gap-3">
            <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
              Save the room as it is
            </span>

            <div className="flex items-center gap-3">
              <SceneLights preview={preview} live className="w-[72px] shrink-0" />

              <Input
                value={name}
                maxLength={80}
                placeholder="Name this scene"
                aria-label="Scene name"
                enterKeyHint="done"
                className="min-w-0 flex-1"
                onChange={(event) => setName(event.target.value)}
                onKeyDown={(event) => {
                  if (event.key === "Enter") submit();
                }}
              />
            </div>

            {taken ? (
              <p className="text-caption text-danger">
                {roomName} already has a scene called that.
              </p>
            ) : null}

            {/* Plain rows: a caption to the panel above, not a table. */}
            <ul className="grid gap-1">
              {devices.map((device) => (
                <li key={device.id} className="flex items-baseline justify-between gap-4">
                  <span className="min-w-0 truncate text-caption text-subtle">{device.name}</span>
                  <span className="shrink-0 font-mono text-caption text-muted tabular">
                    {describe(device, targets)}
                  </span>
                </li>
              ))}
            </ul>

            <Button onClick={submit} disabled={!saveable || save.isPending} className="mt-1">
              {save.isPending ? "Saving" : "Save scene"}
            </Button>
          </section>
        </DrawerBody>
      </DrawerContent>
    </Drawer>
  );
}
