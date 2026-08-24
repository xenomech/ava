import { cn } from "@ava/ui";
import { useRef, type KeyboardEvent, type PointerEvent } from "react";

/** Plate is 288 tall with 12 of padding; the paddle is 130, so it travels 134. */
const TRAVEL = 134;
/** Past this the drag counts as a flick rather than a tap. */
const FLICK = 12;
/** How far the paddle gives under a finger. It rocks, it does not slide. */
const NUDGE = 16;

/**
 * The one control a room has.
 *
 * A room is not dimmable — a plug and a batten share nothing but power — so
 * the only honest gestures are up and down. Both fire unconditionally: flicking
 * up in a room that is already half on means "all of it", not "no change".
 * A tap alternates, for anyone who would rather not gesture at all.
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
  const start = useRef<number | null>(null);

  /* Written straight to the node. A pointermove fires far too often to be
     worth a render, and nothing else on the page reads the offset. */
  const rock = (offset: number) => {
    const node = paddle.current;
    if (!node) return;

    node.style.transitionDuration = "0ms";
    node.style.transform = `translate3d(0, ${(on ? 0 : TRAVEL) + offset}px, 0)`;
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
    start.current = event.clientY;
  };

  const move = (event: PointerEvent<HTMLButtonElement>) => {
    if (start.current === null) return;

    const dy = event.clientY - start.current;
    rock(Math.max(-NUDGE, Math.min(NUDGE, dy)));
  };

  const up = (event: PointerEvent<HTMLButtonElement>) => {
    if (start.current === null) return;

    const dy = event.clientY - start.current;
    start.current = null;
    release();

    if (dy < -FLICK) onFlick(true);
    else if (dy > FLICK) onFlick(false);
    else onFlick(!on);
  };

  const cancel = () => {
    start.current = null;
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
      className="relative grid place-items-center"
      style={{ "--lit": color } as React.CSSProperties}
    >
      <span
        aria-hidden
        className={cn(
          "pointer-events-none absolute size-[300px] rounded-full blur-[64px]",
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
          "relative h-[288px] w-[180px] touch-none select-none rounded-[36px] p-3",
          "border border-white/10 bg-gradient-to-b from-white/8 to-black/40",
          "shadow-[inset_0_1px_0_rgb(255_255_255/0.08),0_24px_60px_-20px_rgb(0_0_0/0.9)]",
          "outline-none focus-visible:ring-2 focus-visible:ring-fg focus-visible:ring-offset-4",
          "focus-visible:ring-offset-bg disabled:cursor-not-allowed disabled:opacity-40",
        )}
      >
        <span
          ref={paddle}
          className={cn(
            "absolute inset-x-3 top-3 block h-[130px] rounded-[26px]",
            "transition-transform duration-[190ms] ease-spring will-change-transform",
            on
              ? "translate-y-0 bg-[var(--lit)] shadow-[0_0_44px_-4px_var(--lit),inset_0_2px_0_rgb(255_255_255/0.5)]"
              : "translate-y-[134px] bg-gradient-to-b from-neutral-800 to-neutral-950 shadow-[inset_0_1px_0_rgb(255_255_255/0.06)]",
          )}
        >
          <span
            aria-hidden
            className={cn(
              "absolute left-1/2 top-1/2 h-[3px] w-7 -translate-x-1/2 -translate-y-1/2 rounded-full",
              on ? "bg-black/25" : "bg-white/15",
            )}
          />
        </span>
      </button>
    </div>
  );
}
