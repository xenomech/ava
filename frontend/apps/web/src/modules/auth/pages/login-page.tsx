import { Button, Field, Input } from "@ava/ui";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Link, useNavigate, useSearch } from "@tanstack/react-router";
import { useState } from "react";

import { isApiError } from "@/config/http/request";
import { login } from "../api";
import { AuthCard } from "../components/auth-card";
import { FormError } from "../components/form-error";
import { deviceName } from "../lib/device-name";
import { SESSION_QUERY_KEY } from "../queries";

export function LoginPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const search = useSearch({ strict: false }) as { redirect?: string };

  const [form, setForm] = useState({ email: "", password: "" });

  const attempt = useMutation({
    mutationFn: () => login(form, deviceName()),
    onSuccess: (result) => {
      queryClient.setQueryData(SESSION_QUERY_KEY, { user: result.user, tenant: result.tenant });
      void navigate({ to: search.redirect ?? "/" });
    },
  });

  const fieldErrors = isApiError(attempt.error) ? attempt.error.details : undefined;

  return (
    <AuthCard
      title="Welcome back"
      description="Sign in to your home."
      footer={
        <>
          New here?{" "}
          <Link to="/auth/register" className="tap font-semibold text-fg hover:underline">
            Set up your home
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
        <Field label="Email" error={fieldErrors?.email}>
          <Input
            type="email"
            autoComplete="email"
            placeholder="you@home.com"
            required
            value={form.email}
            onChange={(event) => setForm((f) => ({ ...f, email: event.target.value }))}
          />
        </Field>

        <Field
          label="Password"
          error={fieldErrors?.password}
          action={
            <Link to="/auth/forgot-password" className="tap text-caption text-muted hover:text-fg">
              Forgot?
            </Link>
          }
        >
          <Input
            type="password"
            autoComplete="current-password"
            placeholder="••••••••"
            required
            value={form.password}
            onChange={(event) => setForm((f) => ({ ...f, password: event.target.value }))}
          />
        </Field>

        {fieldErrors ? null : <FormError error={attempt.error} fallback="Could not sign in" />}

        <Button type="submit" block className="mt-1" loading={attempt.isPending}>
          Continue
        </Button>
      </form>
    </AuthCard>
  );
}
