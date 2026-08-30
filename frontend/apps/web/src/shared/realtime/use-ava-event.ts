import { avaEvent, type AvaEvent } from "@ava/contracts";
import { useEffect } from "react";

import { useAvaSocket } from "./socket";

/** Every validated event off the socket, so parsing lives here rather than in each subscriber. */
export function useAvaEvent(onEvent: (event: AvaEvent) => void) {
  const socket = useAvaSocket();

  useEffect(
    () =>
      socket.subscribe((raw) => {
        const parsed = avaEvent.safeParse(JSON.parse(raw));
        if (parsed.success) onEvent(parsed.data);
      }),
    [socket, onEvent],
  );
}
