import { buttonVariants } from "@ava/ui";
import { useMutation } from "@tanstack/react-query";
import { Link, useSearch } from "@tanstack/react-router";
import { useEffect } from "react";

import { AuthCard } from "../components/auth-card";
import { verifyEmail } from "../api";
import { isApiError } from "@/config/http/request";

export function VerifyEmailPage() {
  const { token } = useSearch({ strict: false }) as { token?: string };

  const attempt = useMutation({ mutationFn: (value: string) => verifyEmail(value) });
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
        description="This verification link is missing its token."
        footer={
          <Link to="/auth/login" className="tap font-semibold text-fg">
            Back to sign in
          </Link>
        }
      />
    );
  }

  if (attempt.isSuccess) {
    return (
      <AuthCard title="Email verified" description="Your account is ready. Sign in to continue.">
        <Link to="/auth/login" className={buttonVariants({ className: "w-full" })}>
          Sign in
        </Link>
      </AuthCard>
    );
  }

  if (attempt.isError) {
    return (
      <AuthCard
        title="Link no longer valid"
        description={
          isApiError(attempt.error)
            ? attempt.error.message
            : "This verification link has expired or was already used."
        }
        footer={
          <Link to="/auth/login" className="tap font-semibold text-fg">
            Back to sign in
          </Link>
        }
      />
    );
  }

  return <AuthCard title="Verifying your email" description="One moment." />;
}
