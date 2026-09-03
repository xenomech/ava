import { defineConfig, minimal2023Preset as preset } from "@vite-pwa/assets-generator/config";

// The mark is white on near-black, so the preset's white fill would frame it in the wrong colour.
const tile = "#0b0b0d";

export default defineConfig({
  headLinkOptions: {
    preset: "2023",
  },
  preset: {
    ...preset,
    // Android crops to a circle, so the letter stays inside the safe zone with room to spare.
    maskable: {
      ...preset.maskable,
      padding: 0.2,
      resizeOptions: { fit: "contain", background: tile },
    },
    // iOS applies its own rounding, so this one runs to the edge instead of sitting in a box.
    apple: { ...preset.apple, padding: 0.08, resizeOptions: { fit: "contain", background: tile } },
  },
  images: ["public/logo.png"],
});
