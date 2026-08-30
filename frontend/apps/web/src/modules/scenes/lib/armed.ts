// Which scene this room's switch is pointed at, kept per device, not per account.
const PREFIX = "ava:v1:armed-scene:";

export function armedScene(roomId: string): string | null {
  try {
    return localStorage.getItem(PREFIX + roomId);
  } catch {
    // Private windows and blocked storage: "all on" is a fine default.
    return null;
  }
}

export function armScene(roomId: string, sceneId: string | null): void {
  try {
    if (sceneId === null) localStorage.removeItem(PREFIX + roomId);
    else localStorage.setItem(PREFIX + roomId, sceneId);
  } catch {
    // The chip still shows as armed for this session; only the memory is lost.
  }
}
