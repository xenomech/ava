import { Button, cn } from "@ava/ui";
import { CheckIcon, CopyIcon } from "lucide-react";
import { useState } from "react";

// The one moment the token exists in readable form, so the copy button is the loudest thing here.
export function TokenReveal({ value, onDone }: { value: string; onDone: () => void }) {
  const [copied, setCopied] = useState(false);

  const copy = () => {
    void navigator.clipboard.writeText(value).then(
      () => {
        setCopied(true);
        window.setTimeout(() => setCopied(false), 2000);
      },
      () => setCopied(false),
    );
  };

  return (
    <div className="grid gap-3 rounded-lg border border-warning/40 bg-surface p-4">
      <div className="grid gap-1">
        <b className="text-body font-semibold">Copy this now</b>
        <p className="text-small text-muted">
          This is the only time the token is shown. If you lose it, delete it and make another.
        </p>
      </div>

      <div className="flex items-center gap-2">
        {/* Selectable and wrapping: a token is long and people do copy it by hand. */}
        <code
          className={cn(
            "min-w-0 flex-1 select-all break-all rounded-md bg-raised p-3",
            "font-mono text-caption",
          )}
        >
          {value}
        </code>

        <Button variant="secondary" size="icon" aria-label="Copy token" onClick={copy}>
          {copied ? <CheckIcon className="size-4" /> : <CopyIcon className="size-4" />}
        </Button>
      </div>

      <Button variant="secondary" className="justify-self-start" onClick={onDone}>
        Done
      </Button>
    </div>
  );
}
