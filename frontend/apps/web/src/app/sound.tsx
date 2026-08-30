import { playSound, unlockSound } from "@ava/ui";
import { useEffect } from "react";

/** Anything you press gets the click, unless it says otherwise. */
const PRESSABLE = 'button, a[href], [role="tab"], [role="menuitem"], [role="option"]';

/**
 * The click every control makes, wired once instead of in every component.
 *
 * A delegated listener rather than a call inside each button: there are dozens
 * of pressable things and more arriving, and a rule that has to be remembered
 * at each call site is a rule that ends up applied to about half of them.
 *
 * Controls that mean something more specific — a switch turning a room on, a
 * device toggling — mark themselves `data-sound="none"` and play their own.
 */
export function SoundEffects() {
  useEffect(() => {
    /* Audio is refused until the page has been interacted with. This runs
       first, in the same gesture, so the very press that unlocks also gets its
       click — which is the moment a browser will allow one. */
    const unlock = () => unlockSound();

    const click = (event: PointerEvent) => {
      const target = event.target as Element | null;
      const hit = target?.closest?.(PRESSABLE);

      if (!hit || hit.closest('[data-sound="none"]')) return;
      if (hit.hasAttribute("disabled") || hit.getAttribute("aria-disabled") === "true") return;

      playSound("click");
    };

    document.addEventListener("pointerdown", unlock, { once: true, capture: true });
    /* pointerdown, not click: the sound belongs to the press, and waiting for
       the release puts it audibly behind the finger. */
    document.addEventListener("pointerdown", click, { capture: true });

    return () => {
      document.removeEventListener("pointerdown", unlock, {
        capture: true,
      } as EventListenerOptions);
      document.removeEventListener("pointerdown", click, { capture: true } as EventListenerOptions);
    };
  }, []);

  return null;
}
