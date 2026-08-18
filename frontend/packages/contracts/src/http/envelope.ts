import { z } from "zod";

import { errorBodySchema, type ErrorBody } from "./error-codes";

export type Success<T> = { success: true; data: T };
export type Failure = { success: false; error: ErrorBody };

export type ApiResponse<T> = Success<T> | Failure;

export const failureSchema = z.object({
  success: z.literal(false),
  error: errorBodySchema,
    data: z.unknown().nullish(),
});

export type FailureEnvelope = z.infer<typeof failureSchema>;

export const successEnvelopeSchema = z.object({
  success: z.literal(true),
  data: z.unknown(),
});
