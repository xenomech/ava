import {
  Button,
  Chip,
  ChipToggle,
  CommandDialog,
  CommandEmpty,
  CommandHint,
  CommandInput,
  CommandItem,
  CommandList,
  Device,
  DeviceHalo,
  Drawer,
  DrawerBody,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
  Field,
  FieldControl,
  Input,
  MarkerSlider,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Slider,
  Switch,
  Tabs,
  TabsList,
  TabsTrigger,
  cn,
  type DeviceKind,
} from "@ava/ui";
import { toast } from "sonner";
import { useLayoutEffect, useRef, useState } from "react";

import { useTheme } from "@/shared/components/theme-provider";

const KINDS: DeviceKind[] = [
  "bulb",
  "tube",
  "strip",
  "lamp",
  "plug",
  "sensor",
  "fan",
  "heater",
  "speaker",
];

const SURFACES = [
  ["bg", "bg-bg"],
  ["surface", "bg-surface"],
  ["raised", "bg-raised"],
  ["border", "bg-border"],
  ["border-strong", "bg-border-strong"],
] as const;

const STATUS = [
  ["accent", "bg-accent"],
  ["off", "bg-off"],
  ["success", "bg-success"],
  ["warning", "bg-warning"],
  ["danger", "bg-danger"],
] as const;

const INK = [
  ["fg", "text-fg"],
  ["muted", "text-muted"],
  ["subtle", "text-subtle"],
] as const;

const TYPE_SCALE = [
  ["hero", "text-hero", "Welcome home"],
  ["display", "text-display", "Living room"],
  ["title", "text-title", "Ceiling"],
  ["lead", "text-lead", "Philips Hue White Ambiance"],
  ["body", "text-body", "Body copy sits at this size."],
  ["small", "text-small", "Secondary detail"],
  ["caption", "text-caption", "CAPTION"],
  ["micro", "text-micro", "MICRO"],
] as const;

const TRACKING = [
  ["tighter", "tracking-tighter"],
  ["tight", "tracking-tight"],
  ["snug", "tracking-snug"],
  ["normal", "tracking-normal"],
  ["caps", "tracking-caps"],
] as const;

const RADII = [
  ["xs", "rounded-xs"],
  ["sm", "rounded-sm"],
  ["md", "rounded-md"],
  ["lg", "rounded-lg"],
  ["xl", "rounded-xl"],
  ["2xl", "rounded-2xl"],
] as const;

export function DesignPage() {
  const { resolvedTheme, setTheme } = useTheme();
  const [level, setLevel] = useState(72);
  const [on, setOn] = useState(true);
  const [palette, setPalette] = useState(false);

  return (
    <div className="min-h-full bg-bg p-8">
      <div className="mx-auto grid max-w-5xl gap-10">
        <header className="flex items-end justify-between gap-4">
          <div>
            <h1 className="text-display font-semibold">Design system</h1>
            <p className="mt-1 text-small text-muted">Every primitive, both themes.</p>
          </div>
          <Tabs value={resolvedTheme ?? "dark"} onValueChange={setTheme}>
            <TabsList>
              <TabsTrigger value="dark">Dark</TabsTrigger>
              <TabsTrigger value="light">Light</TabsTrigger>
            </TabsList>
          </Tabs>
        </header>

        <Section title="Palette">
          <div className="grid gap-4">
            <div className="grid grid-cols-5 gap-3">
              {[...SURFACES, ...STATUS].map(([name, klass]) => (
                <div key={name} className="grid gap-2">
                  <div className={`h-14 rounded-sm border border-border ${klass}`} />
                  <code className="font-mono text-caption text-subtle">{name}</code>
                </div>
              ))}
            </div>
            <div className="flex gap-6">
              {INK.map(([name, klass]) => (
                <span key={name} className={`text-body ${klass}`}>
                  {name}
                </span>
              ))}
            </div>
          </div>
        </Section>

        <Section title="Devices — live at three levels">
          <div className="grid gap-6">
            {[100, 45, 0].map((lv) => (
              <div key={lv} className="grid grid-cols-5 gap-4">
                {KINDS.map((kind) => (
                  <div
                    key={kind}
                    className="relative grid h-32 place-items-center rounded-lg border border-border bg-surface"
                    style={{ "--level": lv, "--lit": "#ffb463" } as React.CSSProperties}
                  >
                    <DeviceHalo className="w-3/4" />
                    <Device kind={kind} level={lv} color="#ffb463" className="h-[86%]" />
                  </div>
                ))}
              </div>
            ))}
          </div>
        </Section>

        <Section title="Buttons">
          <div className="flex flex-wrap items-center gap-3">
            <Button>Primary</Button>
            <Button variant="secondary">Secondary</Button>
            <Button variant="ghost">Ghost</Button>
            <Button variant="quiet">Quiet</Button>
            <Button variant="danger">Danger</Button>
            <Button loading>Loading</Button>
            <Button disabled>Disabled</Button>
            <Button size="sm">Small</Button>
          </div>
        </Section>

        <Section title="Controls">
          <div className="grid max-w-md gap-5">
            <div className="flex items-center justify-between">
              <span className="text-small text-muted">Switch</span>
              <Switch checked={on} onCheckedChange={setOn} aria-label="Demo switch" />
            </div>

            <div className="grid gap-2">
              <div className="flex items-baseline justify-between">
                <span className="text-small text-muted">Slider</span>
                <output className="font-mono text-small tabular">{level}%</output>
              </div>
              <Slider
                value={[level]}
                max={100}
                onValueChange={([v]) => setLevel(v ?? 0)}
                aria-label="Demo slider"
              />
            </div>

            <Field label="Email" hint="We never share this.">
              <Input type="email" placeholder="you@home.com" />
            </Field>

            <Field label="Password" error="That password is too short.">
              <Input type="password" defaultValue="abc" />
            </Field>

            <Field label="Role">
              <Select defaultValue="member">
                <FieldControl>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                </FieldControl>
                <SelectContent>
                  <SelectItem value="admin">Admin</SelectItem>
                  <SelectItem value="member">Member</SelectItem>
                  <SelectItem value="guest">Guest</SelectItem>
                </SelectContent>
              </Select>
            </Field>

            <div className="grid gap-2">
              <span className="text-small text-muted">Slider · marker</span>
              <MarkerSlider
                defaultValue={[30]}
                max={100}
                aria-label="Demo temperature"
                style={{
                  background:
                    "linear-gradient(to right, #ff9233, #ffbe7a, #ffe4c4, #ffffff, #d6e6ff, #a8c8ff)",
                }}
              />
            </div>
          </div>
        </Section>

        <Section title="Chips">
          <div className="flex flex-wrap gap-2">
            <Chip>Neutral</Chip>
            <Chip tone="muted">Muted</Chip>
            <Chip tone="success">On since 07:12</Chip>
            <Chip tone="warning">Open</Chip>
            <Chip tone="danger">Low battery</Chip>
            <ChipToggle defaultPressed>Selectable</ChipToggle>
          </div>
        </Section>

        <Section title="Toasts">
          <div className="flex flex-wrap gap-2">
            <Button variant="secondary" size="sm" onClick={() => toast("Kitchen is all on")}>
              Plain
            </Button>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => toast.success("Reading lamp moved to Bedroom")}
            >
              Success
            </Button>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => toast.warning("2 changed, 1 skipped")}
            >
              Warning
            </Button>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => toast.error("The hub did not accept that")}
            >
              Error
            </Button>
            <Button
              variant="secondary"
              size="sm"
              onClick={() =>
                toast("A new version of Ava is ready", {
                  description: "Reload to pick up the latest build.",
                  action: { label: "Reload", onClick: () => undefined },
                })
              }
            >
              With an action
            </Button>
          </div>
        </Section>

        <Section title="Drawer">
          <Drawer>
            <DrawerTrigger asChild>
              <Button variant="ghost">Open drawer</Button>
            </DrawerTrigger>
            <DrawerContent>
              <DrawerHeader>
                <div>
                  <DrawerTitle>Add a device</DrawerTitle>
                  <DrawerDescription>Pair something new to this home</DrawerDescription>
                </div>
              </DrawerHeader>
              <DrawerBody>
                <div className="grid grid-cols-2 gap-3 pb-4">
                  {KINDS.slice(0, 4).map((kind) => (
                    <div
                      key={kind}
                      className="grid gap-2 rounded-lg border border-border p-3"
                      style={{ "--level": 70, "--lit": "#ffb463" } as React.CSSProperties}
                    >
                      <div className="grid h-16 place-items-center">
                        <Device kind={kind} level={70} color="#ffb463" className="h-full" />
                      </div>
                      <span className="text-small font-semibold capitalize">{kind}</span>
                    </div>
                  ))}
                </div>
              </DrawerBody>
              <DrawerFooter>
                <DrawerClose asChild>
                  <Button variant="ghost">Cancel</Button>
                </DrawerClose>
                <Button className="flex-1">Continue</Button>
              </DrawerFooter>
            </DrawerContent>
          </Drawer>
        </Section>

        <Section title="Command palette">
          <Button variant="ghost" onClick={() => setPalette(true)}>
            Open palette
          </Button>
          <CommandDialog open={palette} onOpenChange={setPalette}>
            <CommandInput placeholder="Search devices and rooms…" />
            <CommandList>
              <CommandEmpty>Nothing matches that.</CommandEmpty>
              <CommandItem>
                Ceiling <CommandHint>Living room</CommandHint>
              </CommandItem>
              <CommandItem>
                Desk lamp <CommandHint>Study</CommandHint>
              </CommandItem>
              <CommandItem>
                Front door sensor <CommandHint>Hallway</CommandHint>
              </CommandItem>
            </CommandList>
          </CommandDialog>
        </Section>

        <Section title="Type scale">
          <div className="grid gap-2">
            {TYPE_SCALE.map(([name, klass, text]) => (
              <TypeRow key={name} name={name} klass={klass} text={text} />
            ))}
          </div>
        </Section>

        <Section title="Tracking">
          <div className="grid gap-2">
            {TRACKING.map(([name, klass]) => (
              <TypeRow
                key={name}
                name={name}
                klass={`text-lead ${klass}`}
                text="Living room ceiling"
              />
            ))}
          </div>
        </Section>

        <Section title="Radius">
          <div className="flex flex-wrap gap-3">
            {RADII.map(([name, klass]) => (
              <div key={name} className="grid gap-2">
                <div className={`size-16 border border-border-strong bg-surface ${klass}`} />
                <code className="font-mono text-caption text-subtle">{name}</code>
              </div>
            ))}
          </div>
        </Section>
      </div>
    </div>
  );
}

function TypeRow({ name, klass, text }: { name: string; klass: string; text: string }) {
  const sample = useRef<HTMLSpanElement>(null);
  const [metrics, setMetrics] = useState("");

  useLayoutEffect(() => {
    if (!sample.current) return;
    const style = getComputedStyle(sample.current);
    setMetrics(`${style.fontSize} / ${style.lineHeight} / ${style.letterSpacing}`);
  }, [klass]);

  return (
    <div className="flex flex-wrap items-baseline gap-x-4 gap-y-1">
      <code className="w-20 shrink-0 font-mono text-caption text-subtle">{name}</code>
      {/* The sample must be allowed to shrink: at a phone width the measured
          size on the right was pushing the row past the viewport, which widens
          the layout viewport and quietly rescales the whole page. */}
      <span ref={sample} className={cn("min-w-0 flex-1 truncate", klass)}>
        {text}
      </span>
      <code className="ml-auto shrink-0 font-mono text-caption text-subtle tabular">{metrics}</code>
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="grid gap-4">
      <h2 className="text-caption font-semibold uppercase tracking-caps text-subtle">{title}</h2>
      {children}
    </section>
  );
}
