import { z } from "zod";
import { sceneTargetDto } from "./scene.dto";

export const createSceneRequest = z.object({
  name: z.string().min(1).max(80),
  targets: z.array(sceneTargetDto).min(1).max(100),
});

export type CreateSceneRequest = z.infer<typeof createSceneRequest>;
