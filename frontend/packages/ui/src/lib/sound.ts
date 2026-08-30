import { SOUNDS, type SoundName } from "./sounds.data";

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
let enabled = read();
let unlocked = false;

const listeners = new Set<(on: boolean) => void>();
const pool = new Map<SoundName, HTMLAudioElement>();

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
  if (typeof Audio === "undefined") return null;

  let audio = pool.get(name);

  if (!audio) {
    audio = new Audio(SOUNDS[name]);
    audio.preload = "auto";
    audio.volume = VOLUME;
    pool.set(name, audio);
  }

  return audio;
}

export function playSound(name: SoundName): void {
  if (!enabled || !unlocked) return;

  const audio = element(name);
  if (!audio) return;

  /* Cloned per play so a rapid second press is not swallowed by the first
     still being in flight. */
  const voice = audio.cloneNode() as HTMLAudioElement;
  voice.volume = VOLUME;
  void voice.play().catch(() => undefined);
}

/**
 * Browsers refuse to play audio until the page has been interacted with, and
 * an attempt before then can leave the element in a wedged state. This marks
 * the moment it becomes allowed; the app calls it from the first real gesture.
 */
export function unlockSound(): void {
  unlocked = true;
}

export function soundEnabled(): boolean {
  return enabled;
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
