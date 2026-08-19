import {
  Button,
  Chip,
  cn,
  Device,
  Drawer,
  DrawerBody,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  Field,
  Input,
  type DeviceKind,
} from "@ava/ui";
import { WifiIcon } from "lucide-react";
import { useEffect, useRef, useState } from "react";

import { Steps } from "@/shared/components/steps";

const KINDS: { kind: DeviceKind; label: string }[] = [
  { kind: "bulb", label: "Light" },
  { kind: "strip", label: "Light strip" },
  { kind: "plug", label: "Smart plug" },
  { kind: "sensor", label: "Sensor" },
];

const STEPS = ["Type", "Pair", "Place", "Done"] as const;

export function AddDeviceDrawer({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const [at, setAt] = useState(0);
  const [kind, setKind] = useState<DeviceKind>("bulb");
  const [name, setName] = useState("Desk lamp");

  const resetTimer = useRef<ReturnType<typeof setTimeout>>(undefined);

  useEffect(() => () => clearTimeout(resetTimer.current), []);

  const close = () => {
    onOpenChange(false);
    clearTimeout(resetTimer.current);
    resetTimer.current = setTimeout(() => setAt(0), 400);
  };

  const advance = (direction: 1 | -1) => {
    const next = at + direction;
    if (next < 0 || next >= STEPS.length) return close();
    setAt(next);
  };

  return (
    <Drawer open={open} onOpenChange={onOpenChange}>
      <DrawerContent>
        <DrawerHeader>
          <div>
            <DrawerTitle>Add a device</DrawerTitle>
            <DrawerDescription>Pair something new to this home</DrawerDescription>
          </div>
        </DrawerHeader>

        <DrawerBody className="grid content-start gap-4 pb-4">
          <Steps labels={STEPS} at={at} />

          {at === 0 ? (
            <div className="grid grid-cols-2 gap-3">
              {KINDS.map((option) => {
                const selected = option.kind === kind;

                return (
                  <button
                    key={option.kind}
                    type="button"
                    aria-pressed={selected}
                    onClick={() => setKind(option.kind)}
                    className={cn(
                      "grid gap-2 rounded-lg border border-border p-3 text-left",
                      "transition-colors duration-150 ease-out hover:border-border-strong",
                      selected && "border-fg",
                    )}
                    style={
                      { "--level": selected ? 75 : 0, "--lit": "#ffb463" } as React.CSSProperties
                    }
                  >
                    <span className="grid h-16 place-items-center">
                      <Device
                        kind={option.kind}
                        level={selected ? 75 : 0}
                        color="#ffb463"
                        className="h-full"
                      />
                    </span>
                    <span className="text-small font-semibold">{option.label}</span>
                  </button>
                );
              })}
            </div>
          ) : null}

          {at === 1 ? (
            <div className="grid justify-items-center gap-4 text-center">
              <p className="text-small text-muted">Switch the light off and on three times.</p>
              <div
                className="grid h-44 w-full place-items-center"
                style={{ "--level": 85, "--lit": "#ffb463" } as React.CSSProperties}
              >
                <Device kind={kind} level={85} color="#ffb463" className="h-full" />
              </div>
              <Chip tone="muted">
                <WifiIcon aria-hidden /> Searching your network…
              </Chip>
            </div>
          ) : null}

          {at === 2 ? (
            <div className="grid gap-4">
              <Field label="Device name">
                {(props) => (
                  <Input {...props} value={name} onChange={(e) => setName(e.target.value)} />
                )}
              </Field>
              <div className="grid gap-2">
                <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
                  Room
                </span>
                <div className="flex flex-wrap gap-2">
                  {["Living room", "Bedroom", "Kitchen", "Studio"].map((room, i) => (
                    <Chip key={room} tone={i === 3 ? "neutral" : "muted"}>
                      {room}
                    </Chip>
                  ))}
                </div>
              </div>
            </div>
          ) : null}

          {at === 3 ? (
            <div className="grid justify-items-center gap-3 text-center">
              <div
                className="grid h-44 w-full place-items-center"
                style={{ "--level": 100, "--lit": "#ffb463" } as React.CSSProperties}
              >
                <Device kind={kind} level={100} color="#ffb463" className="h-full" />
              </div>
              <b className="text-title font-semibold">{name} is ready</b>
              <p className="text-small text-muted">Added to Studio.</p>
            </div>
          ) : null}
        </DrawerBody>

        <DrawerFooter>
          <Button variant="ghost" onClick={() => advance(-1)}>
            Back
          </Button>
          <Button className="flex-1" onClick={() => advance(1)}>
            {at === STEPS.length - 1 ? "Done" : "Continue"}
          </Button>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
}
