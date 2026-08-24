import type { RoomDto } from "@ava/contracts";
import {
  Button,
  Menu,
  MenuContent,
  MenuItem,
  MenuSeparator,
  MenuTrigger,
  cn,
} from "@ava/ui";
import { ChevronDownIcon, ChevronUpIcon, MoreHorizontalIcon, PencilIcon, Trash2Icon } from "lucide-react";
import { useEffect, useRef, useState } from "react";

/**
 * The room's name, and everything you can do to the room itself.
 *
 * Renaming, reordering and deleting used to sit in a row beside the title. On a
 * phone that meant three permanent buttons crowding the heading, with a
 * delete — which takes the room and unassigns everything in it — a thumb's
 * width from the name. They are behind one control now, so the heading is the
 * name and nothing else until you ask for more.
 */
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
    <div className="flex min-w-0 items-center gap-1">
      {/* Tapping the name still renames it — the fastest path stays the
          shortest one, and the menu repeats it for anyone who would not guess. */}
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

      <Menu>
        <MenuTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            aria-label={`Actions for ${room.name}`}
            className="shrink-0 text-muted [@media(hover:hover)]:size-9"
          >
            <MoreHorizontalIcon className="size-4" aria-hidden />
          </Button>
        </MenuTrigger>

        <MenuContent>
          <MenuItem onSelect={() => setEditing(true)}>
            <PencilIcon aria-hidden />
            Rename
          </MenuItem>
          <MenuItem disabled={isFirst} onSelect={() => onMove(-1)}>
            <ChevronUpIcon aria-hidden />
            Move up
          </MenuItem>
          <MenuItem disabled={isLast} onSelect={() => onMove(1)}>
            <ChevronDownIcon aria-hidden />
            Move down
          </MenuItem>

          <MenuSeparator />

          <MenuItem tone="danger" onSelect={onDelete}>
            <Trash2Icon aria-hidden />
            {deviceCount === 0 ? "Delete room" : `Delete, freeing ${deviceCount}`}
          </MenuItem>
        </MenuContent>
      </Menu>
    </div>
  );
}
