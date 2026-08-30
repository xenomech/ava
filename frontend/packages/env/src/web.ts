import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

/**
 * Where the API lives, as either an absolute URL or a same-origin path.
 *
 * A deployed build wants a path: the web server proxies /api to the API service
 * so the browser stays on one origin, which is what keeps the SameSite=Lax
 * session cookies working. Only `z.url()` was accepted before, so a path failed
 * validation and the app threw before it rendered anything.
 *
 * Local development still points at another port, so an absolute URL has to
 * keep working too.
 */
const apiUrl = z
  .string()
  .refine(
    (value) => value.startsWith("/") || URL.canParse(value),
    "must be an absolute URL or a path beginning with /",
  );

export const env = createEnv({
  clientPrefix: "VITE_",
  client: {
    VITE_API_URL: apiUrl.default("http://localhost:8000/api/v1"),
  },
  runtimeEnv: import.meta.env,
  emptyStringAsUndefined: true,
});
