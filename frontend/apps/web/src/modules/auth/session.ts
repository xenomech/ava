import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

import { logout } from "./api";
import { SESSION_QUERY_KEY, sessionQuery } from "./queries";

export function useSession() {
  const { data, isPending } = useQuery(sessionQuery);

  return {
    user: data?.user ?? null,
    tenant: data?.tenant ?? null,
    isAuthenticated: Boolean(data),
    isLoading: isPending,
    isAdmin: data?.tenant.role === "owner" || data?.tenant.role === "admin",
  };
}

export function useSignOut() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: () => logout(),
    onSettled: async () => {
      queryClient.setQueryData(SESSION_QUERY_KEY, null);
      await queryClient.invalidateQueries();
      void navigate({ to: "/auth/login" });
    },
  });
}
