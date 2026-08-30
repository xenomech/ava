import { type FlowStepDto } from "@ava/contracts";

import { BigField } from "./big-field";
import { HubStep } from "./hub-step";

const SUGGESTIONS = ["Home", "The flat", "Rowan Street"];

export function StepFields({
  step,
  values,
  onChange,
}: {
  step: FlowStepDto;
  values: Record<string, string>;
  onChange: (name: string, value: string) => void;
}) {
  const errors = step.errors;

  switch (step.id) {
    case "home":
      return (
        <div className="grid gap-5">
          <BigField
            label="Home name"
            placeholder="Home"
            autoComplete="off"
            required
            value={values.name ?? ""}
            error={errors.name}
            hint="You can change this whenever you like."
            onChange={(event) => onChange("name", event.target.value)}
          />

          <div className="flex flex-wrap gap-2">
            {SUGGESTIONS.map((suggestion) => (
              <button
                key={suggestion}
                type="button"
                onClick={() => onChange("name", suggestion)}
                className="min-h-11 rounded-full border border-border px-4 text-small text-muted transition-colors duration-150 ease-out hover:border-border-strong hover:text-fg"
              >
                {suggestion}
              </button>
            ))}
          </div>
        </div>
      );

    case "hub":
      return (
        <HubStep
          code={values.user_code ?? ""}
          onCodeChange={(value) => onChange("user_code", value)}
          error={errors.user_code}
        />
      );

    default:
      return (
        <p className="text-lead text-muted">
          This step is not supported by this version of the app. Skip it, or update the app.
        </p>
      );
  }
}
