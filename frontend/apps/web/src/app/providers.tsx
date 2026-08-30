import { QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "@ava/ui";
import type { ReactNode } from "react";

import { queryClient } from "./query-client";
import { ThemeProvider } from "@/shared/components/theme-provider";
import { SoundEffects } from "./sound";

export function Providers({ children }: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider
        attribute="data-theme"
        defaultTheme="dark"
        disableTransitionOnChange
        storageKey="ava.theme"
      >
        {children}
        <SoundEffects />
        <Toaster closeButton />
      </ThemeProvider>
    </QueryClientProvider>
  );
}
