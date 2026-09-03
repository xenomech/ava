import { execSync } from "node:child_process";
import { readFileSync } from "node:fs";

import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

const pkg = JSON.parse(readFileSync(new URL("./package.json", import.meta.url), "utf8")) as {
  version: string;
};

function buildVersion() {
  if (process.env.APP_VERSION) return process.env.APP_VERSION;

  try {
    const sha = execSync("git rev-parse --short HEAD", { stdio: ["ignore", "pipe", "ignore"] })
      .toString()
      .trim();

    return `${pkg.version}+${sha}`;
  } catch {
    return `${pkg.version}+dev`;
  }
}

export default defineConfig({
  define: {
    __APP_VERSION__: JSON.stringify(buildVersion()),
  },
  server: {
    port: 3000,
  },
  resolve: {
    tsconfigPaths: true,
  },
  plugins: [
    tailwindcss(),
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
    }),
    react(),
    VitePWA({
      registerType: "prompt",
      manifest: {
        name: "Ava",
        short_name: "Ava",
        description: "Ava - multi-tenant hub control",
        // Installed it is a home-screen app, not a page: no chrome, and no white splash on a cold start.
        display: "standalone",
        orientation: "portrait",
        start_url: "/",
        scope: "/",
        theme_color: "#0b0b0d",
        background_color: "#0b0b0d",
      },
      pwaAssets: { disabled: false, config: true },
      workbox: {
        clientsClaim: true,
        skipWaiting: false,
        cleanupOutdatedCaches: true,
        globPatterns: ["**/*.{js,css,html,ico,png,svg,webmanifest}"],
        // logo.png is the artwork the icons are generated from, not something the app ever loads.
        globIgnores: ["**/logo.png"],
        navigateFallbackDenylist: [/^\/api\//],
      },
      devOptions: { enabled: false },
    }),
  ],
});
