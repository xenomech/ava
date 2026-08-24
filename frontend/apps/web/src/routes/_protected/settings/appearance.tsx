import { cn } from "@ava/ui";
import { createFileRoute } from "@tanstack/react-router";
import { MonitorIcon, MoonIcon, SunIcon } from "lucide-react";

import { Section } from "@/shared/components/page";
import { useTheme } from "@/shared/components/theme-provider";

const CHOICES = [
  { value: "light", label: "Light", icon: SunIcon },
  { value: "dark", label: "Dark", icon: MoonIcon },
  { value: "system", label: "System", icon: MonitorIcon },
] as const;

// The theme toggle used to be an icon in the top bar. With the bar gone it
// belongs here, where it can also offer "follow the system" — which a
// two-state toggle never could.
function Appearance() {
  const { theme, setTheme } = useTheme();

  return (
    <Section title="Theme" description="How Ava looks on this device.">
      <div className="grid gap-2 p-5 sm:grid-cols-3">
        {CHOICES.map(({ value, label, icon: Icon }) => (
          <button
            key={value}
            type="button"
            aria-pressed={theme === value}
            onClick={() => setTheme(value)}
            className={cn(
              "grid min-h-11 place-items-center gap-2 rounded-lg border border-border py-4",
              "text-small font-medium text-muted",
              "transition-colors duration-150 ease-out hover:border-border-strong hover:text-fg",
              "aria-pressed:border-fg aria-pressed:bg-surface aria-pressed:text-fg",
            )}
          >
            <Icon className="size-5" aria-hidden />
            {label}
          </button>
        ))}
      </div>
    </Section>
  );
}

export const Route = createFileRoute("/_protected/settings/appearance")({
  component: Appearance,
});
