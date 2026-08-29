import { useCallback, useSyncExternalStore } from "react";

/* One MediaQueryList per query string for the page's lifetime: getSnapshot runs
   on every render and store check, and window.matchMedia re-parses the query
   each call. Sharing the list also means subscribe and getSnapshot read the
   same object. */
const lists = new Map<string, MediaQueryList>();

function listFor(query: string): MediaQueryList {
  let list = lists.get(query);

  if (!list) {
    list = window.matchMedia(query);
    lists.set(query, list);
  }

  return list;
}

/**
 * Subscribes to a media query.
 *
 * `useSyncExternalStore` rather than an effect: the first render already knows
 * the answer, so a layout that differs between phone and desktop does not have
 * to mount the wrong one and then correct itself.
 */
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
