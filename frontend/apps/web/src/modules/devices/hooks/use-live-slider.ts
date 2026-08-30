import { useCallback, useEffect, useEffectEvent, useRef, useState } from "react";

const LIVE_INTERVAL_MS = 400;

/** A control that owns its value while it is moved, generic so a pad shares the slider clock. */
export function useLiveSlider<T>(
  settled: T,
  preview: (value: T) => void,
  commit: (value: T) => void,
) {
  const [dragging, setDragging] = useState<T | null>(null);

  const timer = useRef<ReturnType<typeof setTimeout>>(undefined);
  const sentAt = useRef(0);

  const sendPreview = useEffectEvent((value: T) => preview(value));
  const sendCommit = useEffectEvent((value: T) => commit(value));

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
        sendPreview(value);

        return;
      }

      cancel();

      timer.current = setTimeout(() => {
        sentAt.current = Date.now();
        sendPreview(value);
      }, LIVE_INTERVAL_MS - elapsed);
    },
    [cancel],
  );

  const release = useCallback(
    (value: T) => {
      cancel();
      sentAt.current = Date.now();
      setDragging(null);
      sendCommit(value);
    },
    [cancel],
  );

  /** Abandon an in-flight drag without committing — an interrupted gesture. */
  const reset = useCallback(() => {
    cancel();
    setDragging(null);
  }, [cancel]);

  return { value: dragging ?? settled, dragging, change, release, reset };
}
