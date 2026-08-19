import { useCallback, useEffect, useRef, type PointerEvent as ReactPointerEvent } from "react";

const AXIS_LOCK = 10;
const HOLD_MS = 420;
const FLICK = 0.4;

export type DeviceGestureHandlers = {
  onTap: () => void;
  onHold: () => void;
  onDim: (delta: number) => void;
  onDimEnd: () => void;
  onSwipe: (direction: 1 | -1) => void;
};

export function useDeviceGesture(handlers: DeviceGestureHandlers) {
  const state = useRef<{
    id: number;
    x: number;
    y: number;
    axis: "x" | "y" | null;
    held: boolean;
    timer: ReturnType<typeof setTimeout>;
    lastX: number;
    lastY: number;
    lastT: number;
    velocity: number;
  } | null>(null);

  useEffect(() => {
    return () => {
      if (state.current) clearTimeout(state.current.timer);
    };
  }, []);

  const onPointerDown = useCallback(
    (event: ReactPointerEvent<HTMLElement>) => {
      const control = (event.target as HTMLElement).closest(
        "button, a, input, select, textarea, [role='slider']",
      );
      if (control && control !== event.currentTarget) return;

      event.currentTarget.setPointerCapture(event.pointerId);

      state.current = {
        id: event.pointerId,
        x: event.clientX,
        y: event.clientY,
        axis: null,
        held: false,
        lastX: event.clientX,
        lastY: event.clientY,
        lastT: performance.now(),
        velocity: 0,
        timer: setTimeout(() => {
          if (!state.current || state.current.axis) return;
          state.current.held = true;
          handlers.onHold();
        }, HOLD_MS),
      };
    },
    [handlers],
  );

  const onPointerMove = useCallback(
    (event: ReactPointerEvent<HTMLElement>) => {
      const s = state.current;
      if (!s || event.pointerId !== s.id) return;

      const dx = event.clientX - s.x;
      const dy = event.clientY - s.y;

      if (!s.axis && Math.hypot(dx, dy) > AXIS_LOCK) {
        s.axis = Math.abs(dx) > Math.abs(dy) ? "x" : "y";
        clearTimeout(s.timer);
      }

      if (s.axis) {
        const now = performance.now();
        const dt = now - s.lastT;
        if (dt > 0) {
          const travelled = s.axis === "x" ? event.clientX - s.lastX : event.clientY - s.lastY;
          s.velocity = travelled / dt;
        }
        s.lastX = event.clientX;
        s.lastY = event.clientY;
        s.lastT = now;
      }

      if (s.axis === "y") {
        const height = event.currentTarget.offsetHeight || 1;
        handlers.onDim((-dy / (height * 0.55)) * 100);
      }
    },
    [handlers],
  );

  const finish = useCallback(
    (event: ReactPointerEvent<HTMLElement>) => {
      const s = state.current;
      if (!s || event.pointerId !== s.id) return;

      clearTimeout(s.timer);
      const dx = event.clientX - s.x;

      if (s.axis === "x") {
        const threshold = event.currentTarget.offsetWidth * 0.22;
        const flicked = Math.abs(dx) > 40 && Math.abs(s.velocity) > FLICK;
        if (dx < -threshold || (flicked && dx < 0)) handlers.onSwipe(1);
        else if (dx > threshold || (flicked && dx > 0)) handlers.onSwipe(-1);
      } else if (s.axis === "y") {
        handlers.onDimEnd();
      } else if (!s.held) {
        handlers.onTap();
      }

      state.current = null;
    },
    [handlers],
  );

  return {
    onPointerDown,
    onPointerMove,
    onPointerUp: finish,
    onPointerCancel: finish,
    style: { touchAction: "none" as const, userSelect: "none" as const },
  };
}
