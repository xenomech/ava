import { useCallback, useEffect, useRef, useState } from "react";

const LIVE_INTERVAL_MS = 400;

export function useLiveSlider(
  settled: number,
  preview: (value: number) => void,
  commit: (value: number) => void,
) {
  const [dragging, setDragging] = useState<number | null>(null);

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
    (value: number) => {
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
    (value: number) => {
      cancel();
      sentAt.current = Date.now();
      setDragging(null);
      commitRef.current(value);
    },
    [cancel],
  );

  return { value: dragging ?? settled, dragging, change, release };
}
