/**
 * Which room this browser was last looking at, so `/` can put someone back
 * where they were instead of somewhere arbitrary.
 *
 * Deliberately not keyed by account. The id is validated against the rooms the
 * current session can actually see before it is used, so a stale id from
 * another tenant simply fails to match and the fallback takes over — which is
 * the same protection a tenant-scoped key would give, without needing the
 * session to be loaded before the redirect can be decided.
 */
const KEY = "ava:v1:last-room";

export function rememberRoom(roomId: string): void {
  try {
    localStorage.setItem(KEY, roomId);
  } catch {
    /* Private windows and blocked storage: the fallback is good enough. */
  }
}

export function recallRoom(): string | null {
  try {
    return localStorage.getItem(KEY);
  } catch {
    return null;
  }
}
