import { CloudOffIcon } from "lucide-react";

export function HubOfflineNotice({ name }: { name: string }) {
  return (
    <output className="flex items-center gap-2.5 border-b border-warning/40 bg-raised px-5 py-2.5 sm:px-6">
      <CloudOffIcon className="size-4 shrink-0 text-warning" aria-hidden />
      <p className="text-small text-muted">
        <b className="font-semibold text-fg">{name}</b> is offline, so its lights cannot be reached.
        Controls return the moment it reconnects.
      </p>
    </output>
  );
}
