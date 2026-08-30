/**
 * Which scene this room's switch is pointed at.
 *
 * Flicking up has to mean something, and once someone has chosen a scene the
 * only sensible thing for it to mean is that scene. Remembering the choice is
 * what makes a scene worth having: you pick "Evening" once and the switch is an
 * evening switch from then on.
 *
 * Kept on the device rather than the account, like the last-visited room. It is
 * a preference about this switch in this hand, not a property of the scene, and
 * two people in a house can reasonably want the switch to mean different
 * things.
 *
 * `null` means the default — everything on.
 */
const PREFIX = "ava:v1:armed-scene:";

export function armedScene(roomId: string): string | null {
  try {
    return localStorage.getItem(PREFIX + roomId);
  } catch {
    /* Private windows and blocked storage: "all on" is a fine default. */
    return null;
  }
}

export function armScene(roomId: string, sceneId: string | null): void {
  try {
    if (sceneId === null) localStorage.removeItem(PREFIX + roomId);
    else localStorage.setItem(PREFIX + roomId, sceneId);
  } catch {
    /* The chip still shows as armed for this session; only the memory is lost. */
  }
}
