// Kept out of the component file so editing a step's markup does not blow away
// Fast Refresh for everything importing this.
export function toStepPayload(stepId: string, values: Record<string, string>): unknown {
  switch (stepId) {
    case "home":
      return { name: values.name ?? "" };

    // Pairing already happened through the hub endpoint; the server reads the
    // hub list rather than anything sent here.
    case "hub":
      return {};

    default:
      return values;
  }
}
