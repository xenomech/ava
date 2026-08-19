import { queryOptions } from "@tanstack/react-query";

import { isApiError } from "@/config/http/request";
import { fetchSession, type Session } from "./api";

export const SESSION_QUERY_KEY = ["session"] as const;

export const sessionQuery = queryOptions<Session | null>({
  queryKey: SESSION_QUERY_KEY,
  staleTime: 60_000,
  retry: false,
  queryFn: async ({ signal }) => {
    try {
      return await fetchSession(signal);
    } catch (error) {
      if (isApiError(error) && (error.status === 401 || error.status === 403)) {
        return null;
      }

      throw error;
    }
  },
});
