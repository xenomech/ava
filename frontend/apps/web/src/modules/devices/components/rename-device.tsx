import type { DeviceDto } from "@ava/contracts";
import { Field, Input } from "@ava/ui";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { updateDevice } from "../api";
import { deviceQueries } from "../queries";

// Devices arrive named by their vendor. The API has always accepted a new name;
// nothing in the app ever offered one, so every fixture kept whatever the bulb
// called itself.
export function RenameDevice({ device }: { device: DeviceDto }) {
  const queryClient = useQueryClient();
  const [draft, setDraft] = useState(device.name);

  // Follow the device when focus moves to another fixture, or when a rename
  // lands from somewhere else. Adjusted during render rather than in an effect:
  // an effect would paint one frame holding the previous fixture's name. The id
  // is part of the key so switching between two devices that happen to share a
  // name still resets.
  const identity = `${device.id}:${device.name}`;
  const [syncedTo, setSyncedTo] = useState(identity);

  if (syncedTo !== identity) {
    setSyncedTo(identity);
    setDraft(device.name);
  }

  const save = useMutation({
    mutationFn: (name: string) => updateDevice(device.id, { name }),
    onSuccess: (updated) => {
      queryClient.setQueryData<DeviceDto[]>(deviceQueries.list().queryKey, (current) =>
        current?.map((entry) => (entry.id === updated.id ? updated : entry)),
      );
    },
    onError: (error) => {
      setDraft(device.name);
      toast.error(isApiError(error) ? error.message : "Could not rename this device");
    },
  });

  const commit = () => {
    const name = draft.trim();

    if (name === "" || name === device.name) {
      setDraft(device.name);

      return;
    }

    save.mutate(name);
  };

  return (
    <Field label="Name">
      {(props) => (
        <Input
          {...props}
          value={draft}
          maxLength={100}
          disabled={save.isPending}
          onChange={(event) => setDraft(event.target.value)}
          onBlur={commit}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              event.preventDefault();
              event.currentTarget.blur();
            }

            if (event.key === "Escape") setDraft(device.name);
          }}
        />
      )}
    </Field>
  );
}
