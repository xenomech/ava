/**
 * Where the phone sheet rests before you drag it.
 *
 * Half the screen, so the device is still visible on the stage above it — the
 * whole point of swapping the room's switch for the device is lost if the sheet
 * covers it. Power and brightness fit in that half; colour and details are the
 * reward for pulling it further. `ROOM_HEIGHT` is the other side of the same
 * number, and the room caps its own height to match.
 */
export const SHEET_SNAP = [0.52, 0.94];
export const ROOM_HEIGHT = "h-[48dvh]";

/** Below this the controls arrive as a sheet; above it, as a column. */
export const BESIDE = "(min-width: 1024px)";
