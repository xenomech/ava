import type { ReactNode } from "react";

import { Reveal } from "@/shared/components/reveal";
import { Room, ROOM_STAGES } from "@/shared/components/room";
import { Wordmark } from "@/shared/components/wordmark";

// Signing in stands in the same room as setup, just a little further along.
export function AuthCard({
  title,
  description,
  children,
  footer,
}: {
  title: string;
  description?: string;
  children?: ReactNode;
  footer?: ReactNode;
}) {
  return (
    <Room light={ROOM_STAGES.waking} className="place-items-center">
      <div className="mx-auto grid w-full max-w-[420px] gap-8">
        <Reveal at={0}>
          <Wordmark />
        </Reveal>

        <header className="grid gap-3">
          <Reveal at={1}>
            <h1 className="text-balance font-semibold tracking-tighter text-[clamp(2rem,6vw,2.75rem)] leading-[1.04]">
              {title}
            </h1>
          </Reveal>

          {description ? (
            <Reveal at={2}>
              <p className="text-pretty text-lead text-muted">{description}</p>
            </Reveal>
          ) : null}
        </header>

        {children ? <Reveal at={3}>{children}</Reveal> : null}

        {footer ? (
          <Reveal at={4}>
            <div className="text-small text-muted">{footer}</div>
          </Reveal>
        ) : null}
      </div>
    </Room>
  );
}
