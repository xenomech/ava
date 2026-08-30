// Kept out of the component file so editing a step's markup keeps Fast Refresh.
export function toStepPayload(stepId: string, values: Record<string, string>): unknown {
  switch (stepId) {
    case "home":
      return { name: values.name ?? "" };

    // Pairing already happened, so the server reads the hub list, not this.
    case "hub":
      return {};

    default:
      return values;
  }
}
