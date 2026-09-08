import { createFileRoute } from "@tanstack/react-router";

import { TokensPage } from "@/modules/tokens/pages/tokens-page";

export const Route = createFileRoute("/_protected/settings/tokens")({
  component: TokensPage,
});
