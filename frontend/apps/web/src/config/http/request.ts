import { successEnvelopeSchema } from "@ava/contracts";
import type { AxiosRequestConfig, Method } from "axios";
import type { z } from "zod";

import { apiClient } from "./client";
import { ApiError } from "./types";

type RequestConfig = {
  url: string;
  method?: Method;
  body?: unknown;
  params?: AxiosRequestConfig["params"];
  headers?: AxiosRequestConfig["headers"];
  signal?: AbortSignal;
  skipAuthRefresh?: boolean;
};

export async function request<Schema extends z.ZodType>(
  config: RequestConfig & { schema: Schema },
): Promise<z.infer<Schema>>;
export async function request(config: RequestConfig & { schema?: undefined }): Promise<void>;
export async function request<Schema extends z.ZodType>({
  url,
  method = "get",
  body,
  schema,
  skipAuthRefresh,
  ...options
}: RequestConfig & { schema?: Schema }): Promise<z.infer<Schema> | void> {
  const response = await apiClient.request<unknown>({
    url,
    method,
    data: body,
    ...options,
    ...(skipAuthRefresh ? { skipAuthRefresh: true } : {}),
  } as AxiosRequestConfig);

  if (!schema) {
    return;
  }

  const envelope = successEnvelopeSchema.parse(response.data);

  return schema.parse(envelope.data);
}

export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError;
}
