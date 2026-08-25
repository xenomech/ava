import {
  TRAIT_BRIGHTNESS,
  TRAIT_COLOR,
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  type DeviceDto,
  type TraitValue,
} from "@ava/contracts";

/**
 * How long a write stays this browser's business before the house wins.
 *
 * Long enough to cover the slow path — a dropped packet to a bulb can put the
 * hub's reply over a second behind, and the broker and socket add their own —
 * short enough that a command which genuinely never landed corrects itself
 * while the person is still looking at the room.
 */
const HOLD_MS = 6_000;

/** `deviceId:trait` → what we asked for, and when. */
const pending = new Map<string, { value: TraitValue; at: number }>();

const keyFor = (deviceID: string, trait: string) => `${deviceID}:${trait}`;

/**
 * Show a change now, before anything has confirmed it.
 *
 * A light that waits for a round trip before it looks switched feels broken
 * even when it is working perfectly, so every write paints first.
 */
export function applyLocally(device: DeviceDto, trait: string, value: TraitValue): DeviceDto {
  const state = { ...device.state, [trait]: value };

  if (trait === TRAIT_BRIGHTNESS && typeof value === "number") {
    state[TRAIT_POWER] = value > 0;
  }

  /* A light shows a colour or a temperature, never both, and the hub clears
     whichever one you did not set. The optimistic copy has to say the same
     thing: while it claimed a light had both, the colour panel read the stale
     temperature, stayed on the White tab, and so refused to switch back to it. */
  if (trait === TRAIT_COLOR) delete state[TRAIT_COLOR_TEMP];
  if (trait === TRAIT_COLOR_TEMP) delete state[TRAIT_COLOR];

  return { ...device, state };
}

/**
 * Say that this browser owns a trait until it hears otherwise.
 *
 * Without this, a change was optimistic only until the next thing the server
 * said — and the server says something every thirty seconds whether or not
 * anything happened. Setting a light to 51% and letting go was enough: the next
 * routine sweep carried the old 100%, replaced the cache wholesale, and put the
 * slider back under the person's hand. The value then jumped a second time when
 * the real state finally arrived. Two corrections, neither of them wanted.
 */
export function claim(deviceID: string, trait: string, value: TraitValue): void {
  pending.set(keyFor(deviceID, trait), { value, at: Date.now() });
}

/** Give a device's writes back, on a rejection or a hard refresh. */
export function release(deviceID: string, trait?: string): void {
  if (trait !== undefined) {
    pending.delete(keyFor(deviceID, trait));

    return;
  }

  for (const key of pending.keys()) {
    if (key.startsWith(`${deviceID}:`)) pending.delete(key);
  }
}

/**
 * A device as the server describes it, with anything still being written here
 * left alone.
 *
 * A claim is dropped as soon as the server agrees with it, so the moment a
 * change is confirmed the house is authoritative again and nothing is held back
 * longer than it has to be.
 */
export function reconcile(incoming: DeviceDto): DeviceDto {
  const now = Date.now();
  let device = incoming;

  for (const [key, held] of pending) {
    const [deviceID, trait] = split(key);
    if (deviceID !== incoming.id) continue;

    if (now - held.at > HOLD_MS) {
      pending.delete(key);

      continue;
    }

    if (incoming.state[trait] === held.value) {
      pending.delete(key);

      continue;
    }

    device = applyLocally(device, trait, held.value);
  }

  return device;
}

/** Every device in a list, reconciled. */
export function reconcileAll(incoming: DeviceDto[]): DeviceDto[] {
  return pending.size === 0 ? incoming : incoming.map(reconcile);
}

function split(key: string): [string, string] {
  const at = key.indexOf(":");

  return [key.slice(0, at), key.slice(at + 1)];
}
