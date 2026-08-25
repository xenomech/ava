import { cn } from "@ava/ui";
import { createFileRoute } from "@tanstack/react-router";
import { MonitorIcon, MoonIcon, SunIcon } from "lucide-react";

import { Page, Section } from "@/shared/components/page";
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
    /* This was the one settings page rendering a bare Section, so it ran edge
       to edge while every sibling sat inside the page gutter. */
    <Page>
      <Section title="Theme" description="How Ava looks on this device.">
        <div className="grid grid-cols-3 gap-2 p-4 sm:p-5">
          {CHOICES.map(({ value, label, icon: Icon }) => (
            <button
              key={value}
              type="button"
              aria-pressed={theme === value}
              onClick={() => setTheme(value)}
              /* Three across rather than stacked. A three-way choice read as a
                 tall column of near-identical cards, which took most of a phone
                 screen to say very little. */
              className={cn(
                "grid min-h-[76px] place-items-center gap-2 rounded-lg border border-border px-2 py-3",
                "text-caption font-medium text-muted sm:text-small",
                "transition-colors duration-150 ease-out hover:border-border-strong hover:text-fg",
                "aria-pressed:border-fg aria-pressed:bg-raised aria-pressed:text-fg",
              )}
            >
              <Icon className="size-5" aria-hidden />
              {label}
            </button>
          ))}
        </div>
      </Section>
    </Page>
  );
}

export const Route = createFileRoute("/_protected/settings/appearance")({
  component: Appearance,
});
