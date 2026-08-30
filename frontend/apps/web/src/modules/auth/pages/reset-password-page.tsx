import { Button, Field, Input } from "@ava/ui";
import { useMutation } from "@tanstack/react-query";
import { Link, useNavigate, useSearch } from "@tanstack/react-router";
import { useState } from "react";

import { AuthCard } from "../components/auth-card";
import { FormError } from "../components/form-error";
import { resetPassword } from "../api";

export function ResetPasswordPage() {
  const navigate = useNavigate();
  const { token } = useSearch({ strict: false }) as { token?: string };
  const [password, setPassword] = useState("");

  const attempt = useMutation({
    mutationFn: () => resetPassword(token ?? "", password),
    onSuccess: () => {
      void navigate({ to: "/auth/login" });
    },
  });

  if (!token) {
    return (
      <AuthCard
        title="Link incomplete"
        description="This reset link is missing its token. Request a new one."
        footer={
          <Link to="/auth/forgot-password" className="tap font-semibold text-fg">
            Request a new link
          </Link>
        }
      />
    );
  }

  return (
    <AuthCard
      title="Choose a new password"
      description="Setting a new password signs you out everywhere."
      footer={
        <Link to="/auth/forgot-password" className="hover:text-fg">
          Request a new link
        </Link>
      }
    >
      <form
        className="grid gap-4"
        onSubmit={(event) => {
          event.preventDefault();
          attempt.mutate();
        }}
      >
        <Field label="New password" hint="At least 8 characters">
          <Input
            type="password"
            autoComplete="new-password"
            required
            minLength={8}
            value={password}
            onChange={(event) => setPassword(event.target.value)}
          />
        </Field>

        <FormError error={attempt.error} fallback="Could not reset the password" />

        <Button type="submit" className="mt-1 w-full" disabled={attempt.isPending}>
          {attempt.isPending ? "Saving…" : "Set new password"}
        </Button>
      </form>
    </AuthCard>
  );
}
