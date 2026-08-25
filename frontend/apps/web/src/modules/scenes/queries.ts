import { queryOptions } from "@tanstack/react-query";

import { listScenes } from "./api";

export const sceneQueries = {
  all: () => ["scenes"] as const,
  room: (roomID: string) =>
    queryOptions({
      queryKey: [...sceneQueries.all(), roomID] as const,
      queryFn: ({ signal }) => listScenes(roomID, signal),
    }),
};
