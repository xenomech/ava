import type { CreateApiTokenRequest } from "@ava/contracts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { isApiError } from "@/config/http/request";
import { createToken, deleteToken, revokeToken } from "../api";
import { tokenQueries } from "../queries";

function reason(error: unknown, fallback: string) {
  return isApiError(error) ? error.message : fallback;
}

export function useTokens() {
  const query = useQuery(tokenQueries.list());

  return { tokens: query.data ?? [], isPending: query.isPending, isError: query.isError };
}

export function useTokenActions() {
  const queryClient = useQueryClient();

  const refresh = () => queryClient.invalidateQueries({ queryKey: tokenQueries.all() });

  const create = useMutation({
    mutationFn: (body: CreateApiTokenRequest) => createToken(body),
    onSuccess: () => void refresh(),
    onError: (error) => toast.error(reason(error, "Could not create the token")),
  });

  const revoke = useMutation({
    mutationFn: revokeToken,
    onSuccess: () => {
      toast.success("Token revoked");
      void refresh();
    },
    onError: (error) => toast.error(reason(error, "Could not revoke the token")),
  });

  const remove = useMutation({
    mutationFn: deleteToken,
    onSuccess: () => {
      toast.success("Token deleted");
      void refresh();
    },
    onError: (error) => toast.error(reason(error, "Could not delete the token")),
  });

  return { create, revoke, remove };
}
