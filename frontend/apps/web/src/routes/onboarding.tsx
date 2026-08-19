import { createFileRoute } from "@tanstack/react-router";

import { requireAuth } from "@/modules/auth";
import { OnboardingPage } from "@/modules/onboarding/pages/onboarding-page";

export const Route = createFileRoute("/onboarding")({
  beforeLoad: requireAuth,
  component: OnboardingPage,
});
