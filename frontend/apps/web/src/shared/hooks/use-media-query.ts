import { useCallback, useSyncExternalStore } from "react";

// One MediaQueryList per query: matchMedia re-parses on each call, and getSnapshot runs often.
const lists = new Map<string, MediaQueryList>();

function listFor(query: string): MediaQueryList {
  let list = lists.get(query);

  if (!list) {
    list = window.matchMedia(query);
    lists.set(query, list);
  }

  return list;
}

/** Subscribes to a media query, so the first render already knows the answer. */
export function useMediaQuery(query: string): boolean {
  const subscribe = useCallback(
    (notify: () => void) => {
      const list = listFor(query);
      list.addEventListener("change", notify);

      return () => list.removeEventListener("change", notify);
    },
    [query],
  );

  return useSyncExternalStore(
    subscribe,
    () => listFor(query).matches,
    () => false,
  );
}
