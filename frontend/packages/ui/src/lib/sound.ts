import type { SoundName } from "./sounds.data";

const KEY = "ava:v1:sound";
const VOLUME = 0.22;

// A deliberately small set of noises, and everything here fails quietly.
let enabled: boolean | null = null;
let unlocked = false;

const listeners = new Set<(on: boolean) => void>();
const pool = new Map<SoundName, HTMLAudioElement>();

// ~12KB of base64 no first paint needs, so it stays out of the initial chunk.
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
    // On by default: nobody should have to find a setting to get sound.
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

  // The first play may beat the data chunk; play it late rather than dropping it.
  if (!sounds) {
    void load().then(() => voice(name));

    return;
  }

  voice(name);
}

function voice(name: SoundName): void {
  const audio = element(name);
  if (!audio) return;

  // Cloned per play so a rapid second press is not swallowed by the first.
  const clone = audio.cloneNode() as HTMLAudioElement;
  clone.volume = VOLUME;
  void clone.play().catch(() => undefined);
}

/** Marks the first real gesture, after which browsers will actually play audio. */
export function unlockSound(): void {
  unlocked = true;
  // The unlocking gesture starts the fetch, so its own click has the audio ready.
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
