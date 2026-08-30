import { useCallback, useSyncExternalStore } from "react";

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
      const list = window.matchMedia(query);
      list.addEventListener("change", notify);

      return () => list.removeEventListener("change", notify);
    },
    [query],
  );

  return useSyncExternalStore(
    subscribe,
    () => window.matchMedia(query).matches,
    () => false,
  );
}
