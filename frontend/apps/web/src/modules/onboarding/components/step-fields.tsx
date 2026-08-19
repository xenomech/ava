import {
  Field,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@ava/ui";
import { type FlowStepDto, type OnboardingMetadataDto } from "@ava/contracts";

export function StepFields({
  step,
  metadata,
  values,
  onChange,
}: {
  step: FlowStepDto;
  metadata: OnboardingMetadataDto | null;
  values: Record<string, string>;
  onChange: (name: string, value: string) => void;
}) {
  const errors = step.errors;

  switch (step.id) {
    case "profile":
      return (
        <>
          <Field label="Your name" error={errors.name}>
            {(props) => (
              <Input
                {...props}
                required
                autoComplete="name"
                value={values.name ?? ""}
                onChange={(event) => onChange("name", event.target.value)}
              />
            )}
          </Field>

          <Field label="Phone" hint="Optional." error={errors.phone}>
            {(props) => (
              <Input
                {...props}
                type="tel"
                autoComplete="tel"
                value={values.phone ?? ""}
                onChange={(event) => onChange("phone", event.target.value)}
              />
            )}
          </Field>
        </>
      );

    case "workspace":
      return (
        <Field label="Home name" error={errors.name}>
          {(props) => (
            <Input
              {...props}
              required
              value={values.name ?? ""}
              onChange={(event) => onChange("name", event.target.value)}
            />
          )}
        </Field>
      );

    case "invite_team":
      return (
        <>
          <Field
            label="Email addresses"
            hint="Separate multiple addresses with commas."
            error={errors.emails}
          >
            {(props) => (
              <Input
                {...props}
                value={values.emails ?? ""}
                placeholder="ana@example.com, sam@example.com"
                onChange={(event) => onChange("emails", event.target.value)}
              />
            )}
          </Field>

          {metadata ? (
            <Field label="Role" error={errors.role}>
              {({ id, invalid, ...props }) => (
                <Select
                  value={values.role ?? metadata.invite_roles[0]?.value ?? ""}
                  onValueChange={(role) => onChange("role", role)}
                >
                  <SelectTrigger id={id} invalid={invalid} {...props}>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {metadata.invite_roles.map((option) => (
                      <SelectItem key={option.value} value={option.value}>
                        {option.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            </Field>
          ) : null}
        </>
      );

    default:
      return (
        <p className="text-small text-muted">
          This step is not supported by this version of the app. Skip it, or update the app.
        </p>
      );
  }
}

export function toStepPayload(stepId: string, values: Record<string, string>): unknown {
  switch (stepId) {
    case "profile":
      return { name: values.name ?? "", phone: values.phone || undefined };

    case "workspace":
      return { name: values.name ?? "" };

    case "invite_team":
      return {
        emails: (values.emails ?? "")
          .split(",")
          .map((entry) => entry.trim())
          .filter(Boolean),
        role: values.role || undefined,
      };

    default:
      return values;
  }
}
