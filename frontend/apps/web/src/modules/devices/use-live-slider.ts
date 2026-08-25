import { useCallback, useEffect, useRef, useState } from "react";

const LIVE_INTERVAL_MS = 400;

/**
 * A control that owns its value while it is being moved.
 *
 * Generic in the value, so a two-dimensional pad throttles on exactly the same
 * clock as a slider rather than growing its own copy of this timing. Two
 * versions of "send at most every 400ms, and always on release" would drift
 * apart the first time either was touched.
 */
export function useLiveSlider<T>(
  settled: T,
  preview: (value: T) => void,
  commit: (value: T) => void,
) {
  const [dragging, setDragging] = useState<T | null>(null);

  const timer = useRef<ReturnType<typeof setTimeout>>(undefined);
  const sentAt = useRef(0);
  const previewRef = useRef(preview);
  const commitRef = useRef(commit);

  useEffect(() => {
    previewRef.current = preview;
    commitRef.current = commit;
  });

  const cancel = useCallback(() => {
    clearTimeout(timer.current);
    timer.current = undefined;
  }, []);

  useEffect(() => cancel, [cancel]);

  const change = useCallback(
    (value: T) => {
      setDragging(value);

      const elapsed = Date.now() - sentAt.current;

      if (elapsed >= LIVE_INTERVAL_MS) {
        sentAt.current = Date.now();
        previewRef.current(value);

        return;
      }

      cancel();

      timer.current = setTimeout(() => {
        sentAt.current = Date.now();
        previewRef.current(value);
      }, LIVE_INTERVAL_MS - elapsed);
    },
    [cancel],
  );

  const release = useCallback(
    (value: T) => {
      cancel();
      sentAt.current = Date.now();
      setDragging(null);
      commitRef.current(value);
    },
    [cancel],
  );

  return { value: dragging ?? settled, dragging, change, release };
}
