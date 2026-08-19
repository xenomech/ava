import { detect } from "detect-browser";

export function deviceName(): string {
  const detected = detect();

  if (!detected) {
    return "Browser";
  }

  const browser = detected.name.charAt(0).toUpperCase() + detected.name.slice(1);
  const platform = detected.os;

  return platform ? `${browser} on ${platform}` : browser;
}
