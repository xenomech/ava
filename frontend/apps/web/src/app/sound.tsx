import { playSound, unlockSound } from "@ava/ui";
import { useEffect } from "react";

/** Anything you press gets the click, unless it says otherwise. */
const PRESSABLE = 'button, a[href], [role="tab"], [role="menuitem"], [role="option"]';

/** The click every control makes, delegated once instead of remembered at each call site. */
export function SoundEffects() {
  useEffect(() => {
    // Audio is refused until first interaction, so the press that unlocks also gets its click.
    const unlock = () => unlockSound();

    const click = (event: PointerEvent) => {
      const target = event.target as Element | null;
      const hit = target?.closest?.(PRESSABLE);

      if (!hit || hit.closest('[data-sound="none"]')) return;
      if (hit.hasAttribute("disabled") || hit.getAttribute("aria-disabled") === "true") return;

      playSound("click");
    };

    document.addEventListener("pointerdown", unlock, { once: true, capture: true });
    // pointerdown, not click: waiting for the release puts the sound audibly behind the finger.
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
