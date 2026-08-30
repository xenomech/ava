import type { SoundName } from "./sounds.data";

const KEY = "ava:v1:sound";
const VOLUME = 0.22;

/**
 * The noises the interface makes.
 *
 * Deliberately small: a click for anything you press, and a rising and falling
 * pair for anything that turns on or off. A distinct sound per interaction gets
 * annoying by the third one, and this is an app people open several times an
 * evening to do one thing.
 *
 * Everything here fails quietly. Audio is blocked until a page has been
 * interacted with, storage throws in private windows, and neither is worth an
 * error in a lighting app.
 */
let enabled: boolean | null = null;
let unlocked = false;

const listeners = new Set<(on: boolean) => void>();
const pool = new Map<SoundName, HTMLAudioElement>();

/* The audio data is ~12KB of base64 that no first paint needs, so it stays out
   of the initial chunk and is fetched when sound first becomes possible. */
let sounds: Record<SoundName, string> | null = null;
let loading: Promise<void> | null = null;

function load(): Promise<void> {
  loading ??= import("./sounds.data").then(
    (module) => {
      sounds = module.SOUNDS;
    },
    () => undefined,
  );

  return loading;
}

function isEnabled(): boolean {
  enabled ??= read();

  return enabled;
}

function read(): boolean {
  try {
    /* On by default: the person asking for sound should not have to find a
       setting to get it. */
    return localStorage.getItem(KEY) !== "off";
  } catch {
    return true;
  }
}

function element(name: SoundName): HTMLAudioElement | null {
  if (typeof Audio === "undefined" || !sounds) return null;

  let audio = pool.get(name);

  if (!audio) {
    audio = new Audio(sounds[name]);
    audio.preload = "auto";
    audio.volume = VOLUME;
    pool.set(name, audio);
  }

  return audio;
}

export function playSound(name: SoundName): void {
  if (!isEnabled() || !unlocked) return;

  /* The very first play may land before the data chunk has arrived; it plays a
     few milliseconds late rather than being dropped — the page already has
     user activation, so a slightly deferred play is still allowed. */
  if (!sounds) {
    void load().then(() => voice(name));

    return;
  }

  voice(name);
}

function voice(name: SoundName): void {
  const audio = element(name);
  if (!audio) return;

  /* Cloned per play so a rapid second press is not swallowed by the first
     still being in flight. */
  const clone = audio.cloneNode() as HTMLAudioElement;
  clone.volume = VOLUME;
  void clone.play().catch(() => undefined);
}

/**
 * Browsers refuse to play audio until the page has been interacted with, and
 * an attempt before then can leave the element in a wedged state. This marks
 * the moment it becomes allowed; the app calls it from the first real gesture.
 */
export function unlockSound(): void {
  unlocked = true;
  /* The unlocking gesture also starts the data fetch, so its own click — and
     everything after it — has the audio ready. */
  void load();
}

export function soundEnabled(): boolean {
  return isEnabled();
}

export function setSoundEnabled(next: boolean): void {
  enabled = next;

  try {
    localStorage.setItem(KEY, next ? "on" : "off");
  } catch {
    /* Nothing to do — the setting simply will not survive a reload. */
  }

  for (const listener of listeners) listener(next);

  if (next) playSound("click");
}

export function onSoundChange(listener: (on: boolean) => void): () => void {
  listeners.add(listener);

  return () => listeners.delete(listener);
}
