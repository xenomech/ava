import type { AxiosInstance, InternalAxiosRequestConfig } from "axios";

import { ApiError } from "./types";

type RetriableConfig = InternalAxiosRequestConfig & {
  retried?: boolean;
  skipAuthRefresh?: boolean;
};

let refreshInFlight: Promise<boolean> | null = null;

function refreshOnce(client: AxiosInstance): Promise<boolean> {
  refreshInFlight ??= client
    .post("/auth/refresh", undefined, { skipAuthRefresh: true } as never)
    .then(() => true)
    .catch(() => false)
    .finally(() => {
      refreshInFlight = null;
    });

  return refreshInFlight;
}

export function attachInterceptors(client: AxiosInstance): void {
  client.interceptors.response.use(
    (response) => response,
    async (error: unknown) => {
      const apiError = ApiError.from(error);
      const config = (error as { config?: RetriableConfig }).config;

      const canRetry =
        apiError.status === 401 && config != null && !config.retried && !config.skipAuthRefresh;

      if (!canRetry) {
        return Promise.reject(apiError);
      }

      config.retried = true;

      if (!(await refreshOnce(client))) {
        return Promise.reject(apiError);
      }

      return client.request(config);
    },
  );
}
