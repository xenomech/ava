import type { SceneDto, SceneTargetDto } from "@ava/contracts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback, useState } from "react";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { useApplyTargets } from "@/modules/devices";
import { createScene, deleteScene } from "./api";
import { armScene, armedScene } from "./armed";
import { sceneQueries } from "./queries";

function reason(error: unknown, fallback: string) {
  return isApiError(error) ? error.message : fallback;
}

export function useScenes(roomId: string) {
  const query = useQuery(sceneQueries.room(roomId));

  return { scenes: query.data ?? [], isPending: query.isPending };
}

/**
 * Which scene the room's switch is pointed at, and how to change it.
 *
 * The stored id is checked against the scenes that still exist rather than
 * trusted, so a scene deleted on another device leaves the switch back at "all
 * on" instead of pointed at nothing.
 */
export function useArmedScene(roomId: string, scenes: SceneDto[]) {
  const [chosen, setChosen] = useState(() => armedScene(roomId));

  const arm = useCallback(
    (sceneId: string | null) => {
      armScene(roomId, sceneId);
      setChosen(sceneId);
    },
    [roomId],
  );

  const armed = scenes.find((scene) => scene.id === chosen) ?? null;

  return { armed, arm };
}

export function useSceneActions(roomId: string) {
  const queryClient = useQueryClient();

  const refresh = () =>
    queryClient.invalidateQueries({ queryKey: sceneQueries.room(roomId).queryKey });

  const save = useMutation({
    mutationFn: ({ name, targets }: { name: string; targets: SceneTargetDto[] }) =>
      createScene(roomId, { name, targets }),
    onSuccess: (created) => {
      queryClient.setQueryData<SceneDto[]>(sceneQueries.room(roomId).queryKey, (current) => [
        ...(current ?? []),
        created,
      ]);
      toast.success(`Saved ${created.name}`);
    },
    onError: (error) => toast.error(reason(error, "Could not save the scene")),
  });

  const remove = useMutation({
    mutationFn: (sceneId: string) => deleteScene(roomId, sceneId),
    onSuccess: (_result, sceneId) => {
      queryClient.setQueryData<SceneDto[]>(sceneQueries.room(roomId).queryKey, (current) =>
        current?.filter((scene) => scene.id !== sceneId),
      );
    },
    onError: (error) => {
      toast.error(reason(error, "Could not delete the scene"));
      void refresh();
    },
  });

  return { save, remove };
}

/**
 * Play a scene back.
 *
 * The saved targets go straight through the ordinary batch apply, so a scene
 * behaves exactly like flicking the switch — same optimistic update, same
 * report of anything the hub refused.
 */
export function useApplyScene() {
  const apply = useApplyTargets();

  return useCallback((scene: SceneDto) => apply(scene.targets), [apply]);
}
