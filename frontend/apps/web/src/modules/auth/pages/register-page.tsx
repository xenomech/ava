import { Button, Field, Input } from "@ava/ui";
import { useMutation } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { useState } from "react";

import { AuthCard } from "../components/auth-card";
import { FormError } from "../components/form-error";
import { register } from "../api";
import { isApiError } from "@/config/http/request";

// Three fields. The server derives the username and the home's URL, and the
// home is named in onboarding — asking for either here would be asking someone
// to invent an identifier before they have seen a single light.
export function RegisterPage() {
  const [form, setForm] = useState({ name: "", email: "", password: "" });

  const attempt = useMutation({ mutationFn: () => register(form) });

  const fieldErrors = isApiError(attempt.error) ? attempt.error.details : undefined;

  if (attempt.isSuccess) {
    return (
      <AuthCard
        title="Check your email"
        description={`We sent a link to ${form.email}. Open it to finish setting up.`}
        footer={
          <Link to="/auth/login" className="font-semibold text-fg hover:underline">
            Back to sign in
          </Link>
        }
      />
    );
  }

  return (
    <AuthCard
      title="Set up your home"
      description="Takes about a minute. You can name your home and pair a hub next."
      footer={
        <>
          Already have an account?{" "}
          <Link to="/auth/login" className="font-semibold text-fg hover:underline">
            Sign in
          </Link>
        </>
      }
    >
      <form
        className="grid gap-4"
        onSubmit={(event) => {
          event.preventDefault();
          attempt.mutate();
        }}
      >
        <Field label="Your name" error={fieldErrors?.name}>
          {(props) => (
            <Input
              {...props}
              autoComplete="name"
              placeholder="Alex Mercer"
              required
              value={form.name}
              onChange={(event) => setForm({ ...form, name: event.target.value })}
            />
          )}
        </Field>

        <Field label="Email" error={fieldErrors?.email}>
          {(props) => (
            <Input
              {...props}
              type="email"
              autoComplete="email"
              placeholder="you@home.com"
              required
              value={form.email}
              onChange={(event) => setForm({ ...form, email: event.target.value })}
            />
          )}
        </Field>

        <Field label="Password" error={fieldErrors?.password} hint="At least 8 characters">
          {(props) => (
            <Input
              {...props}
              type="password"
              autoComplete="new-password"
              placeholder="••••••••"
              required
              minLength={8}
              value={form.password}
              onChange={(event) => setForm({ ...form, password: event.target.value })}
            />
          )}
        </Field>

        {fieldErrors ? null : (
          <FormError error={attempt.error} fallback="Could not create your account" />
        )}

        <Button type="submit" block className="mt-1" loading={attempt.isPending}>
          Create account
        </Button>
      </form>
    </AuthCard>
  );
}
