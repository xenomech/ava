import { createFileRoute } from "@tanstack/react-router";

import { MembersPage } from "@/modules/tenant/pages/members-page";

export const Route = createFileRoute("/_protected/settings/people")({
  component: MembersPage,
});
