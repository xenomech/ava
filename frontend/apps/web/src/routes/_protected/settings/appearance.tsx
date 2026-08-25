import { Switch, cn, onSoundChange, setSoundEnabled, soundEnabled } from "@ava/ui";
import { createFileRoute } from "@tanstack/react-router";
import { useSyncExternalStore } from "react";
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
/**
 * Whether the interface makes a noise.
 *
 * `useSyncExternalStore` because the setting lives outside React — the click
 * handler that plays a sound has no hook to read — and this keeps the toggle
 * honest if it is ever changed from somewhere else.
 */
function useSound(): boolean {
  return useSyncExternalStore(onSoundChange, soundEnabled, () => true);
}

function Appearance() {
  const { theme, setTheme } = useTheme();
  const sound = useSound();

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

      <Section title="Sound" description="Small clicks as you press things.">
        {/* A div, not a label: Radix renders the switch as a button, which is
            not a labelable control, so the association has to come from
            aria-label instead. */}
        <div className="flex min-h-14 items-center justify-between gap-4 p-4 sm:p-5">
          <span className="min-w-0">
            <span className="block text-small font-medium">Interface sounds</span>
            <span className="block text-caption text-muted">Kept on this device.</span>
          </span>
          <Switch
            checked={sound}
            data-sound="none"
            onCheckedChange={setSoundEnabled}
            aria-label="Interface sounds"
          />
        </div>
      </Section>
    </Page>
  );
}

export const Route = createFileRoute("/_protected/settings/appearance")({
  component: Appearance,
});
