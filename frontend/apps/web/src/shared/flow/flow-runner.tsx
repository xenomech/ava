import { Button } from "@ava/ui";
import { flowStateDto, type FlowStateDto, type FlowStepDto } from "@ava/contracts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState, type ReactNode } from "react";

import { isApiError } from "@/config/http/request";
import { Loader } from "@/shared/components/loader";
import { Steps } from "@/shared/components/steps";
import { goBack, skipStep, submitStep } from "./api";
import { flowQueries } from "./queries";

export type FlowFieldsProps = {
  step: FlowStepDto;
  metadata: unknown;
  values: Record<string, string>;
  onChange: (name: string, value: string) => void;
};

export type FlowRunnerProps = {
  flowType: string;
  renderFields: (props: FlowFieldsProps) => ReactNode;
  toPayload: (stepId: string, values: Record<string, string>) => unknown;
  onCompleted: () => void;
  loadingLabel?: string;
  errorTitle?: string;
};

export function FlowRunner({
  flowType,
  renderFields,
  toPayload,
  onCompleted,
  loadingLabel = "Loading",
  errorTitle = "Could not load this flow",
}: FlowRunnerProps) {
  const queryClient = useQueryClient();
  const flowQuery = useQuery(flowQueries.state(flowType));

  const [flow, setFlow] = useState<FlowStateDto | null>(null);
  const [values, setValues] = useState<Record<string, string>>({});

  useEffect(() => {
    if (flowQuery.data) setFlow(flowQuery.data);
  }, [flowQuery.data]);

  const current = flow?.steps.find((step) => step.id === flow.current_step) ?? null;

  useEffect(() => {
    setValues({});
  }, [flow?.current_step]);

  useEffect(() => {
    if (flow?.status === "completed") onCompleted();
  }, [flow?.status, onCompleted]);

  const advance = useMutation({
    mutationFn: (action: "submit" | "skip" | "back") => {
      if (action === "skip") return skipStep(flowType);
      if (action === "back") return goBack(flowType);

      return submitStep(flowType, current?.id ?? "", toPayload(current?.id ?? "", values));
    },
    onSuccess: (next) => {
      setFlow(next);
      queryClient.setQueryData(flowQueries.key(flowType), next);
    },
    onError: (error) => {
      if (isApiError(error) && error.status === 422) {
        const rejected = flowStateDto.safeParse(error.data);
        if (rejected.success) setFlow(rejected.data);
      }
    },
  });

  if (flowQuery.isPending && !flow) return <Loader label={loadingLabel} />;

  if (flowQuery.isError && !flow) {
    return (
      <Shell>
        <div className="grid gap-3">
          <h1 className="text-title font-semibold">{errorTitle}</h1>
          <p className="text-small text-muted">
            {isApiError(flowQuery.error) ? flowQuery.error.message : "Please try again."}
          </p>
          <Button block onClick={() => void flowQuery.refetch()}>
            Retry
          </Button>
        </div>
      </Shell>
    );
  }

  if (!flow || !current) return <Loader label={loadingLabel} />;

  const at = flow.steps.findIndex((step) => step.id === current.id);
  const generalError =
    advance.error && isApiError(advance.error) && advance.error.status !== 422
      ? advance.error.message
      : null;

  return (
    <Shell>
      <Steps labels={flow.steps.map((step) => step.title)} at={Math.max(at, 0)} />

      <header className="mt-6 mb-5 grid gap-1.5">
        <h1 className="text-display font-semibold text-balance">{current.title}</h1>
        {current.description ? (
          <p className="text-small text-muted">{current.description}</p>
        ) : null}
      </header>

      <form
        className="grid gap-4"
        onSubmit={(event) => {
          event.preventDefault();
          advance.mutate("submit");
        }}
      >
        {renderFields({
          step: current,
          metadata: flow.metadata,
          values,
          onChange: (name, value) => setValues((previous) => ({ ...previous, [name]: value })),
        })}

        {generalError ? (
          <p role="alert" className="text-small text-danger">
            {generalError}
          </p>
        ) : null}

        <div className="mt-2 flex items-center gap-2">
          {at > 0 ? (
            <Button type="button" variant="ghost" onClick={() => advance.mutate("back")}>
              Back
            </Button>
          ) : null}

          <div className="flex-1" />

          {current.skippable ? (
            <Button type="button" variant="ghost" onClick={() => advance.mutate("skip")}>
              Skip
            </Button>
          ) : null}

          <Button type="submit" loading={advance.isPending}>
            Continue
          </Button>
        </div>
      </form>
    </Shell>
  );
}

function Shell({ children }: { children: ReactNode }) {
  return (
    <main className="grid min-h-dvh place-items-center bg-bg p-6">
      <div className="w-full max-w-[440px] rounded-xl border border-border bg-surface p-7">
        {children}
      </div>
    </main>
  );
}
