import { clsx, type ClassValue } from "clsx";
import { extendTailwindMerge } from "tailwind-merge";

const FONT_SIZES = [
  "micro",
  "caption",
  "small",
  "body",
  "lead",
  "title",
  "display",
  "hero",
] as const;

const COLORS = [
  "bg",
  "surface",
  "raised",
  "border",
  "border-strong",
  "fg",
  "muted",
  "subtle",
  "accent",
  "accent-fg",
  "off",
  "success",
  "warning",
  "danger",
  "glass",
  "glass-edge",
  "wire",
  "scrim",
] as const;

const TRACKING = ["tighter", "tight", "snug", "normal", "caps"] as const;

const merge = extendTailwindMerge({
  extend: {
    theme: {
      radius: ["xs", "sm", "md", "lg", "xl", "2xl"],
    },
    classGroups: {
      "font-size": [{ text: [...FONT_SIZES] }],
      "text-color": [{ text: [...COLORS] }],
      "bg-color": [{ bg: [...COLORS] }],
      "border-color": [{ border: [...COLORS] }],
      "ring-color": [{ ring: [...COLORS] }],
      tracking: [{ tracking: [...TRACKING] }],
    },
  },
});

export function cn(...inputs: ClassValue[]): string {
  return merge(clsx(inputs));
}
