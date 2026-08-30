import { CloudOffIcon } from "lucide-react";

export function HubOfflineNotice({ name }: { name: string }) {
  // The left inset leaves room for the phone's fixed top-left menu button.
  return (
    <output className="flex items-center gap-2.5 border-b border-warning/40 bg-raised py-2.5 pl-16 pr-5 md:px-5 lg:px-6">
      <CloudOffIcon className="size-4 shrink-0 text-warning" aria-hidden />
      <p className="text-small text-muted">
        <b className="font-semibold text-fg">{name}</b> is offline, so its lights cannot be reached.
        Controls return the moment it reconnects.
      </p>
    </output>
  );
}
