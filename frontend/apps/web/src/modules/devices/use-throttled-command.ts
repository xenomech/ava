import { useCallback, useEffect, useRef } from "react";

const DEFAULT_INTERVAL_MS = 300;

export function useThrottled<T>(send: (value: T) => void, intervalMs = DEFAULT_INTERVAL_MS) {
  const lastSentAt = useRef(0);
  const pending = useRef<T | null>(null);
  const timer = useRef<ReturnType<typeof setTimeout>>(undefined);

  useEffect(() => () => clearTimeout(timer.current), []);

  return useCallback(
    (value: T) => {
      const now = Date.now();
      const elapsed = now - lastSentAt.current;

      if (elapsed >= intervalMs) {
        lastSentAt.current = now;
        send(value);

        return;
      }

      pending.current = value;
      clearTimeout(timer.current);

      timer.current = setTimeout(() => {
        if (pending.current === null) return;

        lastSentAt.current = Date.now();
        send(pending.current);
        pending.current = null;
      }, intervalMs - elapsed);
    },
    [send, intervalMs],
  );
}
