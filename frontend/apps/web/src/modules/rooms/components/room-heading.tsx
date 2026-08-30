import type { RoomDto } from "@ava/contracts";
import { cn } from "@ava/ui";
import { useEffect, useRef, useState } from "react";

/** The room's name, and only that: tap it to rename, in place. */
export function RoomHeading({
  room,
  onRename,
}: {
  room: RoomDto;
  onRename: (name: string) => void;
}) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(room.name);
  const input = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (editing) input.current?.select();
  }, [editing]);

  const commit = () => {
    const name = draft.trim();

    setEditing(false);

    if (name === "" || name === room.name) {
      setDraft(room.name);

      return;
    }

    onRename(name);
  };

  const cancel = () => {
    setDraft(room.name);
    setEditing(false);
  };

  if (editing) {
    return (
      <input
        ref={input}
        value={draft}
        onChange={(event) => setDraft(event.target.value)}
        onBlur={commit}
        onKeyDown={(event) => {
          if (event.key === "Enter") commit();
          if (event.key === "Escape") cancel();
        }}
        aria-label={`Rename ${room.name}`}
        className={cn(
          "min-w-0 max-w-64 rounded-sm border border-fg bg-surface px-1.5 py-0.5",
          "text-title font-semibold outline-none",
        )}
      />
    );
  }

  return (
    <button
      type="button"
      onClick={() => setEditing(true)}
      aria-label={`Rename ${room.name}`}
      className={cn(
        "-mx-1.5 flex min-h-11 min-w-0 items-center rounded-sm px-1.5",
        "text-title font-semibold hover:bg-raised",
        "[@media(hover:hover)]:min-h-0 [@media(hover:hover)]:py-0.5",
      )}
    >
      <span className="truncate">{room.name}</span>
    </button>
  );
}
