// Which room this browser last saw; the id is validated against the session's rooms.
const KEY = "ava:v1:last-room";

export function rememberRoom(roomId: string): void {
  try {
    localStorage.setItem(KEY, roomId);
  } catch {
    // Private windows and blocked storage: the fallback is good enough.
  }
}

export function recallRoom(): string | null {
  try {
    return localStorage.getItem(KEY);
  } catch {
    return null;
  }
}
