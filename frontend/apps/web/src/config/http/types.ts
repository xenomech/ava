import { failureSchema, type ErrorBody } from "@ava/contracts";
import axios from "axios";

export class ApiError extends Error {
  readonly status: number;
  readonly code: string;
  readonly details: Record<string, string> | undefined;
  readonly data: unknown;

  constructor(
    message: string,
    status: number,
    code: string,
    details?: Record<string, string>,
    data?: unknown,
  ) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
    this.details = details;
    this.data = data;
  }

  static fromBody(body: ErrorBody, status: number, data?: unknown): ApiError {
    return new ApiError(body.message, status, body.code, body.details, data);
  }

  static from(error: unknown): ApiError {
    if (error instanceof ApiError) {
      return error;
    }

    if (axios.isAxiosError(error)) {
      const status = error.response?.status ?? 0;
      if (status === 0) {
        return new ApiError("Could not reach the server", 0, "network_error");
      }

      const failure = failureSchema.safeParse(error.response?.data);
      if (failure.success) {
        return ApiError.fromBody(failure.data.error, status, failure.data.data);
      }

      return new ApiError(error.message, status, "internal_error");
    }

    return new ApiError(
      error instanceof Error ? error.message : "Something went wrong",
      0,
      "internal_error",
    );
  }

  get isRetryable(): boolean {
    return this.status === 0 || this.status >= 500;
  }
}
