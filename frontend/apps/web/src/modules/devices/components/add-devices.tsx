import {
  Button,
  Device,
  Drawer,
  DrawerBody,
  DrawerContent,
  DrawerDescription,
  DrawerTitle,
  cn,
} from "@ava/ui";
import type { DeviceDto } from "@ava/contracts";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { CheckIcon, PlusIcon } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { updateDevice } from "../api";
import { deviceQueries } from "../queries";
import { deviceColor, deviceKind, deviceLevel } from "./device-stage";

/**
 * Putting a device into a room, from the room.
 *
 * Until this existed the only way to assign a device was the room picker inside
 * its own control sheet — and that sheet is reached *through* a room. A device
 * belonging to no room was therefore unreachable the moment a single room
 * existed: nothing linked to it, and the empty-state that listed loose devices
 * only appeared when there were no rooms at all.
 *
 * Devices already in another room are offered too, and say where they are.
 * Moving one is the same operation as adopting a loose one, so there is no
 * reason to make it a different screen.
 */
export function AddDevices({
  roomId,
  roomName,
  candidates,
}: {
  roomId: string;
  roomName: string;
  /** Every device not already in this room. */
  candidates: DeviceDto[];
}) {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  const move = useMutation({
    mutationFn: (device: DeviceDto) => updateDevice(device.id, { room_id: roomId }),
    onSuccess: (updated) => {
      queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) =>
        current?.map((entry) => (entry.id === updated.id ? updated : entry)),
      );
      toast.success(`${updated.name} moved to ${roomName}`);
    },
    onError: (error) =>
      toast.error(isApiError(error) ? error.message : "Could not move the device"),
  });

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        aria-label={`Add a device to ${roomName}`}
        className={cn(
          "grid w-[126px] shrink-0 snap-start place-items-center gap-2 rounded-lg p-2.5",
          "border border-dashed border-border-strong text-muted",
          "transition-colors duration-150 ease-out hover:border-fg hover:text-fg",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fg",
        )}
      >
        <span className="grid size-9 place-items-center rounded-full border border-border-strong">
          <PlusIcon className="size-4" aria-hidden />
        </span>
        <span className="text-small font-semibold">Add device</span>
      </button>

      <Drawer open={open} onOpenChange={setOpen}>
        <DrawerContent className="max-h-[80dvh]">
          <DrawerTitle className="px-5 pt-2 text-title font-semibold">
            Add to {roomName}
          </DrawerTitle>
          <DrawerDescription className="px-5 pb-1 text-small text-muted">
            {candidates.length === 0
              ? "Every device you have is already in here."
              : "Pick a device to move into this room."}
          </DrawerDescription>

          <DrawerBody className="grid content-start gap-2 pb-8 pt-3">
            {candidates.map((device) => (
              <Candidate
                key={device.id}
                device={device}
                busy={move.isPending && move.variables?.id === device.id}
                onPick={() => move.mutate(device)}
              />
            ))}

            {candidates.length === 0 ? (
              <Button variant="secondary" onClick={() => setOpen(false)}>
                Close
              </Button>
            ) : null}
          </DrawerBody>
        </DrawerContent>
      </Drawer>
    </>
  );
}

function Candidate({
  device,
  busy,
  onPick,
}: {
  device: DeviceDto;
  busy: boolean;
  onPick: () => void;
}) {
  const level = deviceLevel(device);
  const color = deviceColor(device);

  return (
    <button
      type="button"
      disabled={busy}
      onClick={onPick}
      style={{ "--level": level, "--lit": color } as React.CSSProperties}
      className={cn(
        "flex min-h-14 w-full items-center gap-3 rounded-lg border border-border bg-surface p-2.5",
        "text-left transition-colors duration-150 ease-out",
        "hover:border-border-strong focus-visible:outline-none focus-visible:ring-2",
        "focus-visible:ring-fg disabled:opacity-50",
      )}
    >
      <span className="grid size-10 shrink-0 place-items-center">
        <Device kind={deviceKind(device)} level={level} color={color} className="h-full" />
      </span>

      <span className="min-w-0 flex-1">
        <span className="block truncate text-small font-semibold">{device.name}</span>
        <span className="block truncate font-mono text-caption text-subtle">
          {device.room ? `in ${device.room}` : "not in a room"}
        </span>
      </span>

      {busy ? <CheckIcon className="size-4 shrink-0 text-muted" aria-hidden /> : null}
    </button>
  );
}
