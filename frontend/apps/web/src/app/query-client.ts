import { QueryClient } from "@tanstack/react-query";

import { isApiError } from "@/config/http/request";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: (failureCount, error) => isApiError(error) && error.isRetryable && failureCount < 2,
    },
    mutations: {
      retry: false,
    },
  },
});
