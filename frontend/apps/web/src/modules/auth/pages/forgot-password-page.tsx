import { Button, Field, Input } from "@ava/ui";
import { useMutation } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { useState } from "react";

import { AuthCard } from "../components/auth-card";
import { FormError } from "../components/form-error";
import { forgotPassword } from "../api";

export function ForgotPasswordPage() {
  const [email, setEmail] = useState("");

  const attempt = useMutation({ mutationFn: () => forgotPassword(email) });

  if (attempt.isSuccess) {
    return (
      <AuthCard
        title="Check your email"
        description={`If an account exists for ${email}, a reset link is on its way. It expires in an hour.`}
        footer={
          <Link to="/auth/login" className="tap font-semibold text-fg">
            Back to sign in
          </Link>
        }
      />
    );
  }

  return (
    <AuthCard
      title="Reset your password"
      description="We'll email you a link to choose a new one."
      footer={
        <Link to="/auth/login" className="tap font-semibold text-fg">
          Back to sign in
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
        <Field label="Email">
          <Input
            type="email"
            autoComplete="email"
            placeholder="you@home.com"
            required
            value={email}
            onChange={(event) => setEmail(event.target.value)}
          />
        </Field>

        <FormError error={attempt.error} fallback="Could not send the reset link" />

        <Button type="submit" className="mt-1 w-full" disabled={attempt.isPending}>
          {attempt.isPending ? "Sending…" : "Send reset link"}
        </Button>
      </form>
    </AuthCard>
  );
}
