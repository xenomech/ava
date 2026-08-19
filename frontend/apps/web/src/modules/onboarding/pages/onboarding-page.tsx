import { onboardingMetadataDto } from "@ava/contracts";
import { useNavigate } from "@tanstack/react-router";
import { useCallback } from "react";

import { FlowRunner } from "@/shared/flow";
import { ONBOARDING_FLOW } from "../flow";
import { StepFields, toStepPayload } from "../components/step-fields";

export function OnboardingPage() {
  const navigate = useNavigate();

  const onCompleted = useCallback(() => {
    void navigate({ to: "/", replace: true });
  }, [navigate]);

  return (
    <FlowRunner
      flowType={ONBOARDING_FLOW}
      loadingLabel="Loading setup"
      errorTitle="Could not start setup"
      toPayload={toStepPayload}
      onCompleted={onCompleted}
      renderFields={({ step, metadata, values, onChange }) => (
        <StepFields
          step={step}
          metadata={onboardingMetadataDto.safeParse(metadata).data ?? null}
          values={values}
          onChange={onChange}
        />
      )}
    />
  );
}
