import { Button, Input } from "@ava/ui";
import { PlusIcon } from "lucide-react";
import { useEffect, useRef, useState } from "react";

export function NewRoom({ onCreate, busy }: { onCreate: (name: string) => void; busy: boolean }) {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const input = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (open) input.current?.focus();
  }, [open]);

  const close = () => {
    setName("");
    setOpen(false);
  };

  const submit = () => {
    const trimmed = name.trim();
    if (trimmed === "") return close();

    onCreate(trimmed);
    close();
  };

  if (!open) {
    return (
      <Button variant="secondary" size="sm" onClick={() => setOpen(true)}>
        <PlusIcon className="size-4" aria-hidden />
        New room
      </Button>
    );
  }

  return (
    <form
      className="flex items-center gap-2"
      onSubmit={(event) => {
        event.preventDefault();
        submit();
      }}
    >
      <Input
        ref={input}
        value={name}
        onChange={(event) => setName(event.target.value)}
        onKeyDown={(event) => {
          if (event.key === "Escape") close();
        }}
        onBlur={submit}
        placeholder="Kitchen"
        maxLength={80}
        aria-label="New room name"
        className="h-9 w-44"
      />
      <Button type="submit" size="sm" disabled={busy || name.trim() === ""}>
        Add
      </Button>
    </form>
  );
}
