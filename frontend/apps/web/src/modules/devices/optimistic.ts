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

/** deviceId → trait → what we asked for, and when. Keyed by device first so
    reconciling one device never has to walk every other device's claims. */
const pending = new Map<string, Map<string, { value: TraitValue; at: number }>>();

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
  /* Writing one of these retires the other, so holding on to the one you have
     just replaced is holding on to something that is no longer true. Both were
     kept independently, and reconcile then applied them in turn — so setting a
     colour and then a temperature left the panel showing whichever claim
     happened to be applied last, while the bulb showed the one actually sent.
     The light was right and the app argued with it for six seconds. */
  if (trait === TRAIT_COLOR) release(deviceID, TRAIT_COLOR_TEMP);
  if (trait === TRAIT_COLOR_TEMP) release(deviceID, TRAIT_COLOR);

  let traits = pending.get(deviceID);

  if (!traits) {
    traits = new Map();
    pending.set(deviceID, traits);
  }

  traits.set(trait, { value, at: Date.now() });
}

/** Give a device's writes back, on a rejection or a hard refresh. */
export function release(deviceID: string, trait?: string): void {
  if (trait === undefined) {
    pending.delete(deviceID);

    return;
  }

  const traits = pending.get(deviceID);
  if (!traits) return;

  traits.delete(trait);
  if (traits.size === 0) pending.delete(deviceID);
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
  const traits = pending.get(incoming.id);
  if (!traits) return incoming;

  const now = Date.now();
  let device = incoming;

  for (const [trait, held] of traits) {
    if (now - held.at > HOLD_MS || incoming.state[trait] === held.value) {
      traits.delete(trait);

      continue;
    }

    device = applyLocally(device, trait, held.value);
  }

  if (traits.size === 0) pending.delete(incoming.id);

  return device;
}

/** Every device in a list, reconciled. */
export function reconcileAll(incoming: DeviceDto[]): DeviceDto[] {
  /* Expire here as well as per device: a device that has stopped reporting
     never reaches `reconcile`, so its claims would otherwise sit in the map —
     and keep the fast path below disabled — for the rest of the session. */
  const now = Date.now();

  for (const [deviceID, traits] of pending) {
    for (const [trait, held] of traits) {
      if (now - held.at > HOLD_MS) traits.delete(trait);
    }

    if (traits.size === 0) pending.delete(deviceID);
  }

  return pending.size === 0 ? incoming : incoming.map(reconcile);
}
