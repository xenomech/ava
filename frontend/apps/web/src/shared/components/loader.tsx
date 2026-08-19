import { Loader2Icon } from "lucide-react";

export function Loader({ label = "Loading" }: { label?: string }) {
  return (
    <output className="flex min-h-dvh items-center justify-center" aria-live="polite">
      <Loader2Icon className="text-muted size-5 animate-spin" aria-hidden="true" />
      <span className="sr-only">{label}</span>
    </output>
  );
}
