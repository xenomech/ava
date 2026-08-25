import { sceneDto, type CreateSceneRequest, type SceneDto } from "@ava/contracts";
import { z } from "zod";

import { request } from "@/config/http/request";

export function listScenes(roomID: string, signal?: AbortSignal): Promise<SceneDto[]> {
  return request({ url: `/rooms/${roomID}/scenes`, schema: z.array(sceneDto), signal });
}

export function createScene(roomID: string, body: CreateSceneRequest): Promise<SceneDto> {
  return request({ url: `/rooms/${roomID}/scenes`, method: "post", body, schema: sceneDto });
}

export function deleteScene(roomID: string, sceneID: string): Promise<void> {
  return request({ url: `/rooms/${roomID}/scenes/${sceneID}`, method: "delete" });
}
