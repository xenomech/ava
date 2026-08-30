import { cn } from "@ava/ui";
import { useRef, type ClipboardEvent, type KeyboardEvent } from "react";

// Boxes show the four-hyphen-four shape, so a wrong-length paste is obvious.
const GROUP = 4;
const LENGTH = GROUP * 2;

// Hoisted: clean runs on every keystroke and render.
const NOT_CODE = /[^A-Z0-9]/g;

const clean = (value: string) => value.toUpperCase().replace(NOT_CODE, "").slice(0, LENGTH);

/** "BXS8568F" -> "BXS8-568F", the shape the server issued. */
const format = (raw: string) =>
  raw.length > GROUP ? `${raw.slice(0, GROUP)}-${raw.slice(GROUP)}` : raw;

export function CodeField({
  label,
  value,
  onChange,
  onComplete,
  error,
  disabled,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  /** Fired once the last character lands, so a full code pairs on its own. */
  onComplete?: (value: string) => void;
  error?: string;
  disabled?: boolean;
}) {
  const cells = useRef<(HTMLInputElement | null)[]>([]);
  const raw = clean(value);

  const focusCell = (index: number) => {
    const target = cells.current[Math.min(Math.max(index, 0), LENGTH - 1)];
    target?.focus();
    target?.select();
  };

  const commit = (next: string, caret: number) => {
    onChange(format(next));
    focusCell(caret);

    if (next.length === LENGTH) onComplete?.(format(next));
  };

  const onCellChange = (index: number, incoming: string) => {
    const typed = clean(incoming);
    if (!typed) return;

    // Typing over a cell replaces from that point, so overtyping behaves as it looks.
    const next = clean(raw.slice(0, index) + typed);

    commit(next, index + typed.length);
  };

  const onCellKeyDown = (index: number, event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Backspace") {
      event.preventDefault();

      // Backspace clears this cell, or steps back if it is already empty.
      const at = raw[index] ? index : index - 1;
      if (at < 0) return;

      commit(clean(raw.slice(0, at) + raw.slice(at + 1)), at);

      return;
    }

    if (event.key === "ArrowLeft") {
      event.preventDefault();
      focusCell(index - 1);
    }

    if (event.key === "ArrowRight") {
      event.preventDefault();
      focusCell(index + 1);
    }
  };

  const onPaste = (event: ClipboardEvent<HTMLInputElement>) => {
    event.preventDefault();

    const pasted = clean(event.clipboardData.getData("text"));
    if (pasted) commit(pasted, pasted.length);
  };

  return (
    <fieldset className="grid min-w-0 gap-3">
      <legend className="mb-3 text-caption font-semibold uppercase tracking-caps text-subtle">
        {label}
      </legend>

      <div className="flex min-w-0 items-center gap-1.5 sm:gap-2">
        {Array.from({ length: LENGTH }, (_, index) => (
          <Cell
            key={index}
            index={index}
            char={raw[index] ?? ""}
            invalid={Boolean(error)}
            disabled={disabled}
            cells={cells}
            onPaste={onPaste}
            onChange={onCellChange}
            onKeyDown={onCellKeyDown}
          />
        ))}
      </div>

      {error ? (
        <p role="alert" className="text-small text-danger">
          {error}
        </p>
      ) : null}
    </fieldset>
  );
}

function Cell({
  index,
  char,
  invalid,
  disabled,
  cells,
  onPaste,
  onChange,
  onKeyDown,
}: {
  index: number;
  char: string;
  invalid: boolean;
  disabled?: boolean;
  cells: React.RefObject<(HTMLInputElement | null)[]>;
  onPaste: (event: ClipboardEvent<HTMLInputElement>) => void;
  onChange: (index: number, value: string) => void;
  onKeyDown: (index: number, event: KeyboardEvent<HTMLInputElement>) => void;
}) {
  const separator = index === GROUP;

  return (
    <>
      {separator ? <span aria-hidden className="h-px w-3 shrink-0 bg-border-strong" /> : null}

      <input
        ref={(node) => {
          cells.current[index] = node;
        }}
        // The first cell carries the autofill hint for the whole group.
        autoComplete={index === 0 ? "one-time-code" : "off"}
        aria-label={`Character ${index + 1} of ${LENGTH}`}
        aria-invalid={invalid || undefined}
        inputMode="text"
        autoCapitalize="characters"
        spellCheck={false}
        disabled={disabled}
        value={char}
        onChange={(event) => onChange(index, event.target.value)}
        onKeyDown={(event) => onKeyDown(index, event)}
        onPaste={onPaste}
        onFocus={(event) => event.currentTarget.select()}
        className={cn(
          "h-14 min-w-0 flex-1 rounded-md border bg-surface text-center",
          // Capped so eight boxes do not stretch into slabs on a wide screen.
          "max-w-[52px]",
          // 16px+ so iOS never zooms; scales up where there is room.
          "font-mono text-[clamp(1rem,4.5vw,1.375rem)] font-semibold uppercase",
          "transition-[border-color,box-shadow] duration-150 ease-out",
          "outline-none focus:border-[var(--lit)]",
          "disabled:cursor-not-allowed disabled:opacity-50",
          invalid ? "border-danger" : "border-border",
        )}
        style={char ? { boxShadow: "inset 0 0 24px -14px var(--lit)" } : undefined}
      />
    </>
  );
}
