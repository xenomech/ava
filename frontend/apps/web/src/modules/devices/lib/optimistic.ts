import {
  TRAIT_BRIGHTNESS,
  TRAIT_COLOR,
  TRAIT_COLOR_TEMP,
  TRAIT_POWER,
  type DeviceDto,
  type TraitValue,
} from "@ava/contracts";

// How long a write stays this browser's business before the house wins.
const HOLD_MS = 6_000;

// deviceId to trait to what we asked for, keyed by device so reconciling one skips the rest.
const pending = new Map<string, Map<string, { value: TraitValue; at: number }>>();

/** Show a change now, because a light that waits for a round trip feels broken. */
export function applyLocally(device: DeviceDto, trait: string, value: TraitValue): DeviceDto {
  const state = { ...device.state, [trait]: value };

  if (trait === TRAIT_BRIGHTNESS && typeof value === "number") {
    state[TRAIT_POWER] = value > 0;
  }

  // A light shows a colour or a temperature, never both, so the local copy must agree.
  if (trait === TRAIT_COLOR) delete state[TRAIT_COLOR_TEMP];
  if (trait === TRAIT_COLOR_TEMP) delete state[TRAIT_COLOR];

  return { ...device, state };
}

/** Say this browser owns a trait, so a routine sweep cannot undo a change under your hand. */
export function claim(deviceID: string, trait: string, value: TraitValue): void {
  // Writing one of these retires the other, so the replaced claim is no longer true.
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

/** The server's device with local writes left alone, each claim dropped once the house agrees. */
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
  // Expire here too: a device that stopped reporting never reaches reconcile to be cleared.
  const now = Date.now();

  for (const [deviceID, traits] of pending) {
    for (const [trait, held] of traits) {
      if (now - held.at > HOLD_MS) traits.delete(trait);
    }

    if (traits.size === 0) pending.delete(deviceID);
  }

  return pending.size === 0 ? incoming : incoming.map(reconcile);
}
