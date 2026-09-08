import { type TokenScope } from "@ava/contracts";

// Grouped by the thing they act on, so the form reads as "what may this touch" rather than a flat list.
export const SCOPE_GROUPS: { label: string; read: TokenScope; write: TokenScope }[] = [
  { label: "Devices", read: "devices:read", write: "devices:write" },
  { label: "Rooms", read: "rooms:read", write: "rooms:write" },
  { label: "Scenes", read: "scenes:read", write: "scenes:write" },
  { label: "Hubs", read: "hubs:read", write: "hubs:write" },
];

/** "devices:write" as "Devices · change". */
export function scopeLabel(scope: TokenScope): string {
  const [subject, access] = scope.split(":");
  const group = SCOPE_GROUPS.find((entry) => entry.read.startsWith(`${subject}:`));

  return `${group?.label ?? subject} · ${access === "write" ? "change" : "view"}`;
}
