import { isApiError } from "@/config/http/request";

export function FormError({ error, fallback }: { error: unknown; fallback: string }) {
  if (!error) {
    return null;
  }

  return (
    <p className="text-danger text-small" role="alert">
      {isApiError(error) ? error.message : fallback}
    </p>
  );
}
