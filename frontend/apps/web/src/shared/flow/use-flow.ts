import { flowStateDto, type FlowStateDto, type FlowStepDto } from "@ava/contracts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";

import { isApiError } from "@/config/http/request";
import { goBack, skipStep, submitStep } from "./api";
import { flowQueries } from "./queries";

export type FlowAction = "submit" | "skip" | "back";

export type Flow = {
  flow: FlowStateDto | null;
  step: FlowStepDto | null;
  /** Index of the current step, or -1 before anything has loaded. */
  at: number;
  values: Record<string, string>;
  setValue: (name: string, value: string) => void;
  advance: (action: FlowAction) => void;
  isAdvancing: boolean;
  /** Errors that are not per-field; field errors ride along on the step. */
  error: string | null;
  isPending: boolean;
  isError: boolean;
  loadError: string | null;
  retry: () => void;
};

/**
 * The flow state machine, with no opinion about how it looks. Onboarding wraps
 * this in a full-screen cinematic layout; anything else can render it as a
 * plain form. Keeping the machine separate is what lets the two differ without
 * either one growing a `variant` prop.
 */
export function useFlow(
  flowType: string,
  toPayload: (stepId: string, values: Record<string, string>) => unknown,
  onCompleted: () => void,
): Flow {
  const queryClient = useQueryClient();
  const flowQuery = useQuery(flowQueries.state(flowType));

  // The query cache is the only copy of the flow. Mirroring it into component
  // state meant a rejected step wrote to the mirror and left the cache holding
  // a different answer.
  const flow = flowQuery.data ?? null;
  const setFlow = (next: FlowStateDto) =>
    queryClient.setQueryData(flowQueries.key(flowType), next);

  const [values, setValues] = useState<Record<string, string>>({});

  const step = flow?.steps.find((entry) => entry.id === flow.current_step) ?? null;

  // Derived during render rather than in an effect: as an effect this chained
  // off the query and let a new step paint holding the old step's answers.
  const [valuesFor, setValuesFor] = useState(flow?.current_step);

  if (valuesFor !== flow?.current_step) {
    setValuesFor(flow?.current_step);
    setValues({});
  }

  useEffect(() => {
    if (flow?.status === "completed") onCompleted();
  }, [flow?.status, onCompleted]);

  const mutation = useMutation({
    mutationFn: (action: FlowAction) => {
      if (action === "skip") return skipStep(flowType);
      if (action === "back") return goBack(flowType);

      return submitStep(flowType, step?.id ?? "", toPayload(step?.id ?? "", values));
    },
    onSuccess: (next) => setFlow(next),
    onError: (error) => {
      // A 422 carries the rejected flow, field errors and all.
      if (isApiError(error) && error.status === 422) {
        const rejected = flowStateDto.safeParse(error.data);
        if (rejected.success) setFlow(rejected.data);
      }
    },
  });

  return {
    flow,
    step,
    at: flow ? flow.steps.findIndex((entry) => entry.id === flow.current_step) : -1,
    values,
    setValue: (name, value) => setValues((previous) => ({ ...previous, [name]: value })),
    advance: mutation.mutate,
    isAdvancing: mutation.isPending,
    error:
      mutation.error && isApiError(mutation.error) && mutation.error.status !== 422
        ? mutation.error.message
        : null,
    isPending: flowQuery.isPending && !flow,
    isError: flowQuery.isError && !flow,
    loadError: isApiError(flowQuery.error) ? flowQuery.error.message : null,
    retry: () => void flowQuery.refetch(),
  };
}
