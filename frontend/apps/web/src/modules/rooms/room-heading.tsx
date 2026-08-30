import type { RoomDto } from "@ava/contracts";
import { Button, cn } from "@ava/ui";
import { ChevronDownIcon, ChevronUpIcon, Trash2Icon } from "lucide-react";
import { useEffect, useRef, useState } from "react";

export function RoomHeading({
  room,
  deviceCount,
  isFirst,
  isLast,
  onRename,
  onMove,
  onDelete,
}: {
  room: RoomDto;
  deviceCount: number;
  isFirst: boolean;
  isLast: boolean;
  onRename: (name: string) => void;
  onMove: (direction: -1 | 1) => void;
  onDelete: () => void;
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

  return (
    <div className="group/room flex items-center gap-2">
      {editing ? (
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
          className="min-w-0 max-w-64 rounded-sm border border-fg bg-surface px-1.5 py-0.5 text-title font-semibold outline-none"
        />
      ) : (
        <button
          type="button"
          onClick={() => setEditing(true)}
          className="rounded-sm px-1.5 py-0.5 text-title font-semibold -mx-1.5 hover:bg-raised"
          aria-label={`Rename ${room.name}`}
        >
          {room.name}
        </button>
      )}

      <div
        className={cn(
          "flex items-center gap-0.5 transition-opacity",
          "[@media(hover:hover)]:opacity-0",
          "group-hover/room:opacity-100 focus-within:opacity-100",
          editing && "opacity-100",
        )}
      >
        <Button
          variant="ghost"
          size="icon-sm"
          disabled={isFirst}
          onClick={() => onMove(-1)}
          aria-label={`Move ${room.name} up`}
        >
          <ChevronUpIcon className="size-4" aria-hidden />
        </Button>
        <Button
          variant="ghost"
          size="icon-sm"
          disabled={isLast}
          onClick={() => onMove(1)}
          aria-label={`Move ${room.name} down`}
        >
          <ChevronDownIcon className="size-4" aria-hidden />
        </Button>
        <Button
          variant="ghost"
          size="icon-sm"
          onClick={onDelete}
          aria-label={
            deviceCount === 0
              ? `Delete ${room.name}`
              : `Delete ${room.name} and unassign ${deviceCount} devices`
          }
        >
          <Trash2Icon className="size-4" aria-hidden />
        </Button>
      </div>
    </div>
  );
}
