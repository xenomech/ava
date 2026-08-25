import { cn } from "@ava/ui";
import { useRef, type KeyboardEvent, type PointerEvent } from "react";

/**
 * Under this much movement the gesture was a tap, not a drag.
 *
 * A thumb rarely lands still — a deliberate tap on a phone drifts several
 * pixels. The old threshold was 12px in either direction, so an ordinary tap
 * regularly registered as a flick, and which way it went depended on which way
 * the thumb happened to slide.
 */
const TAP_SLOP = 8;

/**
 * A throw this fast commits in its own direction regardless of how far it got.
 * Pixels per millisecond, measured over the last move of the drag.
 */
const FLING = 0.45;

/**
 * The one control a room has.
 *
 * A room is not dimmable — a plug and a batten share nothing but power — so
 * the only honest gestures are up and down. Both fire unconditionally: flicking
 * up in a room that is already half on means "all of it", not "no change".
 * A tap alternates, for anyone who would rather not gesture at all.
 *
 * Three layers, not two. The plate is the housing, the rail is a recess sunk
 * into it, and the paddle travels inside the rail and is clipped by it. Without
 * that middle layer the paddle floats on the front of the plate and the whole
 * thing reads as a picture of a switch rather than a mechanism with a moving
 * part.
 *
 * Everything is sized from `--plate-h`, which shrinks on a short screen. A
 * paddle that is exactly half the rail less its clearance travels exactly its
 * own height, so the resting positions are `top-1.5` and `translate-y-full` at
 * any size and there is no magic number to keep in step with the CSS.
 */
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
  const drag = useRef<{ startY: number; lastY: number; lastAt: number; speed: number } | null>(
    null,
  );

  /* Written straight to the node. A pointermove fires far too often to be worth
     a render, and nothing else on the page reads the offset. */
  const place = (offset: number) => {
    const node = paddle.current;
    if (!node) return;

    node.style.transitionDuration = "0ms";
    node.style.transform = `translate3d(0, ${offset}px, 0)`;
  };

  const release = () => {
    const node = paddle.current;
    if (!node) return;

    node.style.transitionDuration = "";
    node.style.transform = "";
  };

  const down = (event: PointerEvent<HTMLButtonElement>) => {
    if (disabled) return;

    event.currentTarget.setPointerCapture(event.pointerId);
    drag.current = { startY: event.clientY, lastY: event.clientY, lastAt: event.timeStamp, speed: 0 };
  };

  const move = (event: PointerEvent<HTMLButtonElement>) => {
    const state = drag.current;
    const node = paddle.current;
    if (!state || !node) return;

    const elapsed = event.timeStamp - state.lastAt;

    if (elapsed > 0) {
      state.speed = (event.clientY - state.lastY) / elapsed;
      state.lastY = event.clientY;
      state.lastAt = event.timeStamp;
    }

    /* The paddle goes where the finger goes, across its whole travel. It used
       to give by 16px and no more, so a drag the length of the plate moved it
       almost not at all and the control felt inert until you let go. */
    const travel = node.offsetHeight;
    const from = on ? 0 : travel;

    place(Math.min(Math.max(from + (event.clientY - state.startY), 0), travel));
  };

  const up = (event: PointerEvent<HTMLButtonElement>) => {
    const state = drag.current;
    const node = paddle.current;
    if (!state) return;

    drag.current = null;
    release();

    const dy = event.clientY - state.startY;

    if (Math.abs(dy) < TAP_SLOP) {
      onFlick(!on);

      return;
    }

    /* A fast throw wins on its own; a slow drag has to actually carry the
       paddle past the midpoint. Deciding on distance alone made a quick,
       confident flick feel like it had been ignored. */
    if (Math.abs(state.speed) >= FLING) {
      onFlick(state.speed < 0);

      return;
    }

    const travel = node?.offsetHeight ?? 0;
    const landed = (on ? 0 : travel) + dy;

    onFlick(landed < travel / 2);
  };

  const cancel = () => {
    drag.current = null;
    release();
  };

  const key = (event: KeyboardEvent<HTMLButtonElement>) => {
    if (disabled) return;

    if (event.key === "ArrowUp") onFlick(true);
    else if (event.key === "ArrowDown") onFlick(false);
    else if (event.key === " " || event.key === "Enter") onFlick(!on);
    else return;

    event.preventDefault();
  };

  return (
    <div
      className={cn(
        "relative grid place-items-center",
        /* A phone in portrait has to fit a header, this, a readout and a
           carousel. On anything short the plate gives up height first. */
        "[--plate-h:288px] [@media(max-height:760px)]:[--plate-h:240px]",
        "[@media(max-height:640px)]:[--plate-h:200px]",
      )}
      style={{ "--lit": color } as React.CSSProperties}
    >
      {/* The room's light, thrown from behind the plate. The paddle's own glow
          is clipped by the rail, which is right — a recessed lamp spills around
          its housing, not through it. */}
      <span
        aria-hidden
        className={cn(
          "pointer-events-none absolute size-[calc(var(--plate-h)*1.05)] rounded-full blur-[64px]",
          "bg-[var(--lit)] transition-opacity duration-[320ms] ease-out",
          on && !disabled ? "opacity-25" : "opacity-0",
        )}
      />

      <button
        type="button"
        role="switch"
        aria-checked={on}
        aria-label={label}
        disabled={disabled}
        onPointerDown={down}
        onPointerMove={move}
        onPointerUp={up}
        onPointerCancel={cancel}
        onKeyDown={key}
        className={cn(
          "relative h-[var(--plate-h)] w-[calc(var(--plate-h)*0.625)] p-[max(10px,calc(var(--plate-h)*0.049))]",
          "cursor-grab touch-none select-none rounded-[calc(var(--plate-h)*0.14)]",
          "border border-border-strong bg-gradient-to-b from-raised to-surface",
          "shadow-[0_18px_44px_-12px_rgb(0_0_0/0.6),inset_0_1px_0_rgb(255_255_255/0.06)]",
          "outline-none focus-visible:ring-2 focus-visible:ring-fg focus-visible:ring-offset-4",
          "focus-visible:ring-offset-bg active:cursor-grabbing",
          "disabled:cursor-not-allowed disabled:opacity-40",
        )}
      >
        <span
          className={cn(
            "relative block size-full overflow-hidden rounded-[calc(var(--plate-h)*0.105)] bg-bg",
            "shadow-[inset_0_2px_10px_rgb(0_0_0/0.8)]",
          )}
        >
          <span
            ref={paddle}
            className={cn(
              "absolute inset-x-1.5 top-1.5 grid h-[calc(50%-6px)] place-items-center",
              "rounded-[calc(var(--plate-h)*0.077)]",
              "transition-transform duration-[190ms] ease-spring will-change-transform",
              on
                ? "translate-y-0 bg-[var(--lit)] shadow-[0_0_26px_-2px_var(--lit),inset_0_1px_0_rgb(255_255_255/0.4)]"
                : cn(
                    "translate-y-full bg-gradient-to-b from-border-strong to-raised",
                    "shadow-[0_4px_12px_rgb(0_0_0/0.5),inset_0_1px_0_rgb(255_255_255/0.08)]",
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
