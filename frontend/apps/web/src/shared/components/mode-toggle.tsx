import { Button } from "@ava/ui";
import { MoonIcon, SunIcon } from "lucide-react";

import { useTheme } from "./theme-provider";

export function ModeToggle() {
  const { resolvedTheme, setTheme } = useTheme();
  const next = resolvedTheme === "light" ? "dark" : "light";

  return (
    <Button
      variant="ghost"
      size="icon-sm"
      onClick={() => setTheme(next)}
      aria-label={`Switch to ${next} theme`}
    >
      {resolvedTheme === "light" ? <MoonIcon aria-hidden /> : <SunIcon aria-hidden />}
    </Button>
  );
}
