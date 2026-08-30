import { buttonVariants } from "@ava/ui";
import { useMutation } from "@tanstack/react-query";
import { Link, useSearch } from "@tanstack/react-router";
import { useEffect } from "react";

import { AuthCard } from "../components/auth-card";
import { acceptInvite } from "../api";
import { isApiError } from "@/config/http/request";

export function AcceptInvitePage() {
  const { token } = useSearch({ strict: false }) as { token?: string };

  const attempt = useMutation({ mutationFn: (value: string) => acceptInvite(value) });
  const { mutate } = attempt;

  useEffect(() => {
    if (token) {
      mutate(token);
    }
  }, [token, mutate]);

  if (!token) {
    return (
      <AuthCard
        title="Link incomplete"
        description="This invitation link is missing its token."
        footer={
          <Link to="/auth/login" className="font-semibold text-fg">
            Back to sign in
          </Link>
        }
      />
    );
  }

  if (attempt.isSuccess) {
    return (
      <AuthCard
        title="Invitation accepted"
        description="You now have access to that home. Sign in to continue."
      >
        <Link to="/auth/login" className={buttonVariants({ className: "w-full" })}>
          Sign in
        </Link>
      </AuthCard>
    );
  }

  if (attempt.isError) {
    return (
      <AuthCard
        title="Invitation not valid"
        description={
          isApiError(attempt.error)
            ? attempt.error.message
            : "This invitation has expired or was already used."
        }
        footer={
          <Link to="/auth/login" className="font-semibold text-fg">
            Back to sign in
          </Link>
        }
      />
    );
  }

  return <AuthCard title="Accepting your invitation" description="One moment." />;
}
