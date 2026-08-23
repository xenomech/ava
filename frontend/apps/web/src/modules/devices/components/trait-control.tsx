import { Slider, Switch, cn } from "@ava/ui";
import { rangeOf, traitLabel, type CapabilityDto, type TraitValue } from "@ava/contracts";

export function TraitControl({
  capability,
  value,
  disabled,
  onChange,
}: {
  capability: CapabilityDto;
  value: TraitValue | undefined;
  disabled?: boolean;
  onChange: (value: TraitValue) => void;
}) {
  const label = traitLabel(capability.trait);

  if (capability.kind === "bool") {
    return (
      <div className="flex items-center justify-between">
        <Label>{label}</Label>
        <Switch
          checked={value === true}
          disabled={disabled}
          onCheckedChange={onChange}
          aria-label={label}
        />
      </div>
    );
  }

  if (capability.kind === "enum") {
    return (
      <div className="grid gap-2.5">
        <Label>{label}</Label>
        <div className="flex flex-wrap gap-1.5">
          {(capability.values ?? []).map((option) => (
            <button
              key={option}
              type="button"
              disabled={disabled}
              onClick={() => onChange(option)}
              aria-pressed={value === option}
              className={cn(
                "rounded-sm border border-border px-2.5 py-1 text-small capitalize",
                "transition-colors disabled:opacity-40",
                value === option ? "border-fg bg-raised text-fg" : "text-subtle hover:text-fg",
              )}
            >
              {option}
            </button>
          ))}
        </div>
      </div>
    );
  }

  const range = rangeOf(capability);
  if (!range) return null;

  const current = typeof value === "number" ? value : range.min;

  return (
    <div className="grid gap-2.5">
      <div className="flex items-baseline justify-between">
        <Label>{label}</Label>
        <output className="font-mono text-small tabular">
          {current}
          {range.unit}
        </output>
      </div>
      <Slider
        value={[current]}
        min={range.min}
        max={range.max}
        step={range.step}
        disabled={disabled}
        aria-label={label}
        onValueCommit={([next]) => onChange(next ?? range.min)}
      />
    </div>
  );
}

export function TraitReading({
  capability,
  value,
}: {
  capability: CapabilityDto;
  value: TraitValue | undefined;
}) {
  if (value === undefined) return null;

  const shown = typeof value === "boolean" ? (value ? "Yes" : "No") : String(value);

  return (
    <div className="flex items-baseline justify-between">
      <Label>{traitLabel(capability.trait)}</Label>
      <output className="font-mono text-small tabular">
        {shown}
        {capability.unit ?? ""}
      </output>
    </div>
  );
}

function Label({ children }: { children: React.ReactNode }) {
  return (
    <span className="text-caption font-semibold uppercase tracking-caps text-subtle">
      {children}
    </span>
  );
}
