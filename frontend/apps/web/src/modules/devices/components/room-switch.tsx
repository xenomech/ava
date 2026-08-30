import { cn, playSound } from "@ava/ui";
import { useRef, type KeyboardEvent, type PointerEvent } from "react";

// Under this much movement the gesture was a tap, not a drag: a thumb rarely lands still.
const TAP_SLOP = 8;

// A throw this fast commits in its own direction, in pixels per millisecond.
const FLING = 0.45;

/** The one control a room has: up, down, or a tap to alternate. */
export function RoomSwitch({
  on,
  disabled,
  color,
  label,
  onFlick,
}: {
  on: boolean;
  disabled: boolean;
  /** The room's own light, so the paddle glows at the temperature it will be. */
  color: string;
  label: string;
  onFlick: (on: boolean) => void;
}) {
  const paddle = useRef<HTMLSpanElement>(null);

  // The room's own sound, rising or falling, instead of the generic press click.
  const flick = (next: boolean) => {
    playSound(next ? "on" : "off");
    onFlick(next);
  };

  const drag = useRef<{
    pointerId: number;
    startY: number;
    lastY: number;
    lastAt: number;
    speed: number;
  } | null>(null);

  // Written straight to the node, and as `translate` so it replaces the resting utility.
  const place = (offset: number) => {
    const node = paddle.current;
    if (!node) return;

    node.style.transitionDuration = "0ms";
    node.style.translate = `0 ${offset}px`;
  };

  const release = () => {
    const node = paddle.current;
    if (!node) return;

    node.style.transitionDuration = "";
    node.style.translate = "";
  };

  const down = (event: PointerEvent<HTMLButtonElement>) => {
    if (disabled) return;

    // Ignore a second finger: only the pointer that started the drag drives the paddle.
    if (drag.current) return;

    event.currentTarget.setPointerCapture(event.pointerId);
    drag.current = {
      pointerId: event.pointerId,
      startY: event.clientY,
      lastY: event.clientY,
      lastAt: event.timeStamp,
      speed: 0,
    };
  };

  const move = (event: PointerEvent<HTMLButtonElement>) => {
    const state = drag.current;
    const node = paddle.current;
    if (!state || !node || event.pointerId !== state.pointerId) return;

    const elapsed = event.timeStamp - state.lastAt;

    if (elapsed > 0) {
      state.speed = (event.clientY - state.lastY) / elapsed;
      state.lastY = event.clientY;
      state.lastAt = event.timeStamp;
    }

    // The paddle goes where the finger goes, across its whole travel, or it feels inert.
    const travel = node.offsetHeight;
    const from = on ? 0 : travel;

    place(Math.min(Math.max(from + (event.clientY - state.startY), 0), travel));
  };

  const up = (event: PointerEvent<HTMLButtonElement>) => {
    const state = drag.current;
    const node = paddle.current;
    if (!state || event.pointerId !== state.pointerId) return;

    drag.current = null;
    release();

    const dy = event.clientY - state.startY;

    if (Math.abs(dy) < TAP_SLOP) {
      flick(!on);

      return;
    }

    // A fast throw wins on its own, so a quick flick is never mistaken for no change.
    if (Math.abs(state.speed) >= FLING) {
      flick(state.speed < 0);

      return;
    }

    const travel = node?.offsetHeight ?? 0;
    const landed = (on ? 0 : travel) + dy;

    flick(landed < travel / 2);
  };

  const cancel = () => {
    drag.current = null;
    release();
  };

  const key = (event: KeyboardEvent<HTMLButtonElement>) => {
    if (disabled) return;

    if (event.key === "ArrowUp") flick(true);
    else if (event.key === "ArrowDown") flick(false);
    else if (event.key === " " || event.key === "Enter") flick(!on);
    else return;

    event.preventDefault();
  };

  return (
    <div
      className={cn(
        "relative grid place-items-center",
        // One clamp, not stacked breakpoints, so the plate has one answer at every height.
        "[--plate-h:clamp(140px,calc(100dvh-464px),288px)]",
      )}
      style={{ "--lit": color } as React.CSSProperties}
    >
      {/* The room's light, thrown from behind the plate as a recessed lamp would spill. */}
      <span
        aria-hidden
        className={cn(
          "pointer-events-none absolute size-[calc(var(--plate-h)*1.17)] rounded-full",
          "[filter:blur(calc(var(--plate-h)*0.26))]",
          "bg-[var(--lit)] transition-opacity duration-[320ms] ease-out",
          on && !disabled ? "opacity-[0.28]" : "opacity-0",
        )}
      />

      <button
        type="button"
        role="switch"
        aria-checked={on}
        aria-label={label}
        disabled={disabled}
        data-sound="none"
        onPointerDown={down}
        onPointerMove={move}
        onPointerUp={up}
        onPointerCancel={cancel}
        onKeyDown={key}
        className={cn(
          "relative h-[var(--plate-h)] w-[calc(var(--plate-h)*0.625)] p-[max(10px,calc(var(--plate-h)*0.049))]",
          "cursor-grab touch-none select-none rounded-[calc(var(--plate-h)*0.14)]",
          "border border-border-strong bg-gradient-to-b from-raised to-surface switch-plate",
          "outline-none focus-visible:ring-2 focus-visible:ring-fg focus-visible:ring-offset-4",
          "focus-visible:ring-offset-bg active:cursor-grabbing",
          "disabled:cursor-not-allowed disabled:opacity-40",
        )}
      >
        <span
          className={cn(
            "relative block size-full overflow-hidden rounded-[calc(var(--plate-h)*0.105)] bg-bg",
            "switch-rail",
          )}
        >
          <span
            ref={paddle}
            className={cn(
              "absolute inset-x-1.5 top-1.5 grid h-[calc(50%-6px)] place-items-center",
              "rounded-[calc(var(--plate-h)*0.077)]",
              // Position springs; colour and glow follow slowly, so the paddle arrives first.
              "transition-[translate,background-color,box-shadow] will-change-transform",
              "duration-[190ms,300ms,300ms]",
              "[transition-timing-function:cubic-bezier(0.22,1.2,0.36,1),ease-out,ease-out]",
              on
                ? "translate-y-0 bg-[var(--lit)] switch-paddle-on"
                : cn(
                    "translate-y-full bg-gradient-to-b from-border-strong to-raised",
                    "switch-paddle-off",
                  ),
            )}
          >
            <span
              aria-hidden
              className={cn(
                "h-[3px] w-9 rounded-full transition-colors duration-[400ms] ease-out",
                on ? "bg-black/35" : "bg-white/20",
              )}
            />
          </span>
        </span>
      </button>
    </div>
  );
}
