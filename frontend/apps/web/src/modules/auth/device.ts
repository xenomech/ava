import { detect } from "detect-browser";

/* The user agent cannot change for the life of the page, so it is parsed at
   most once however many login attempts there are. */
let name: string | null = null;

export function deviceName(): string {
  name ??= describe();

  return name;
}

function describe(): string {
  const detected = detect();

  if (!detected) {
    return "Browser";
  }

  const browser = detected.name.charAt(0).toUpperCase() + detected.name.slice(1);
  const platform = detected.os;

  return platform ? `${browser} on ${platform}` : browser;
}
