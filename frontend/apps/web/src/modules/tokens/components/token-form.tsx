import type { CreateApiTokenRequest, TokenScope } from "@ava/contracts";
import { Button, Field, Input, Switch, cn } from "@ava/ui";
import { useState } from "react";

import { SCOPE_GROUPS } from "../lib/scopes";

const NEVER = "never";

// Offered rather than free text: an expiry nobody picks is an expiry nobody sets.
const LIFETIMES = [
  { value: "30", label: "30 days" },
  { value: "90", label: "90 days" },
  { value: "365", label: "A year" },
  { value: NEVER, label: "Never" },
] as const;

export function TokenForm({
  busy,
  onCreate,
}: {
  busy: boolean;
  onCreate: (body: CreateApiTokenRequest) => void;
}) {
  const [name, setName] = useState("");
  const [scopes, setScopes] = useState<TokenScope[]>(["devices:read", "devices:write"]);
  const [lifetime, setLifetime] = useState<string>("90");

  const toggle = (scope: TokenScope, on: boolean) =>
    setScopes((current) => (on ? [...current, scope] : current.filter((entry) => entry !== scope)));

  const submit = (event: React.FormEvent) => {
    event.preventDefault();

    onCreate({
      name: name.trim(),
      scopes,
      ...(lifetime === NEVER ? {} : { expires_in_days: Number(lifetime) }),
    });
  };

  const ready = name.trim().length > 0 && scopes.length > 0;

  return (
    <form className="grid gap-5 p-5" onSubmit={submit}>
      <Field label="What is it for" hint="A name you will recognise later, like “Siri shortcuts”.">
        <Input
          required
          maxLength={80}
          placeholder="Siri shortcuts"
          value={name}
          onChange={(event) => setName(event.target.value)}
        />
      </Field>

      <fieldset className="grid gap-2.5">
        <legend className="pb-1 text-caption font-semibold uppercase tracking-caps text-subtle">
          What it may do
        </legend>

        {SCOPE_GROUPS.map((group) => (
          <div key={group.label} className="flex items-center justify-between gap-4">
            <span className="text-small">{group.label}</span>

            <div className="flex shrink-0 items-center gap-4">
              <ScopeSwitch
                label={`View ${group.label.toLowerCase()}`}
                caption="View"
                checked={scopes.includes(group.read)}
                onChange={(on) => toggle(group.read, on)}
              />
              <ScopeSwitch
                label={`Change ${group.label.toLowerCase()}`}
                caption="Change"
                checked={scopes.includes(group.write)}
                onChange={(on) => toggle(group.write, on)}
              />
            </div>
          </div>
        ))}

        <p className="text-caption text-subtle">
          Grant the least that works. A shortcut that turns lights on needs devices only.
        </p>
      </fieldset>

      <fieldset className="grid gap-2">
        <legend className="pb-1 text-caption font-semibold uppercase tracking-caps text-subtle">
          Expires
        </legend>

        <div className="flex flex-wrap gap-1.5">
          {LIFETIMES.map((option) => (
            <button
              key={option.value}
              type="button"
              aria-pressed={lifetime === option.value}
              onClick={() => setLifetime(option.value)}
              className={cn(
                "min-h-9 rounded-full border border-border px-3.5 text-small text-muted",
                "transition-colors duration-150 ease-out hover:text-fg",
                "aria-pressed:border-fg aria-pressed:text-fg",
              )}
            >
              {option.label}
            </button>
          ))}
        </div>
      </fieldset>

      <Button type="submit" className="justify-self-start" disabled={!ready} loading={busy}>
        Create token
      </Button>
    </form>
  );
}

function ScopeSwitch({
  label,
  caption,
  checked,
  onChange,
}: {
  label: string;
  caption: string;
  checked: boolean;
  onChange: (on: boolean) => void;
}) {
  return (
    <label className="flex items-center gap-2">
      <span className="text-caption text-subtle">{caption}</span>
      <Switch checked={checked} onCheckedChange={onChange} aria-label={label} />
    </label>
  );
}
