import { Button, Device } from "@ava/ui";

import { Reveal } from "@/shared/components/reveal";
import { Wordmark } from "@/shared/components/wordmark";
import { kelvinToCss } from "@/shared/lib/kelvin";

// Device writes whatever it is given straight into --lit, so it has to be a
// real colour. Passing "var(--lit)" makes the variable reference itself, which
// resolves to nothing and leaves the bulb dark.
const EMBER = kelvinToCss(1900);

// Before any question, one screen that is just the product's idea: a single
// light in a dark room. Nothing to fill in, one thing to press.
export function Welcome({ onBegin }: { onBegin: () => void }) {
  return (
    <div className="grid justify-items-center gap-10 text-center">
      <Reveal at={0}>
        <Wordmark />
      </Reveal>

      <Reveal at={1}>
        <div className="grid h-[240px] place-items-center">
          <Device kind="bulb" level={55} color={EMBER} className="h-full" />
        </div>
      </Reveal>

      <div className="grid gap-4">
        <Reveal at={3}>
          <h1 className="text-balance font-semibold tracking-tighter text-[clamp(2.5rem,8vw,4.25rem)] leading-[1]">
            Let there be light.
          </h1>
        </Reveal>

        <Reveal at={4}>
          <p className="mx-auto max-w-[38ch] text-pretty text-lead text-muted">
            Every lamp in the house, on one screen. Two questions and you are done.
          </p>
        </Reveal>
      </div>

      <Reveal at={5}>
        <Button size="md" className="min-w-[200px]" onClick={onBegin}>
          Begin
        </Button>
      </Reveal>
    </div>
  );
}
